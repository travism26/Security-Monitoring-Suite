import mongoose from "mongoose";
import { Team, TeamDoc } from "../models/team";
import { User } from "../models/user";
import { NotFoundError, BadRequestError, ForbiddenError } from "../errors";

// Interface for team hierarchy
interface TeamHierarchy extends Omit<TeamDoc, keyof mongoose.Document> {
  children: TeamHierarchy[];
}

export class TeamService {
  // Create a new team
  static async createTeam(
    name: string,
    tenantId: string,
    ownerId: string,
    description?: string,
    parentTeamId?: string,
    settings?: {
      resourceQuota?: {
        maxMembers?: number;
        maxStorage?: number;
        maxProjects?: number;
      };
      visibility?: string;
      allowResourceSharing?: boolean;
    }
  ): Promise<TeamDoc> {
    const session = await mongoose.startSession();
    session.startTransaction();

    try {
      // Verify owner exists and belongs to tenant
      const owner = await User.findOne({
        _id: ownerId,
        tenantId: new mongoose.Types.ObjectId(tenantId),
      });

      if (!owner) {
        throw new BadRequestError("Invalid owner for team");
      }

      // Verify parent team if provided
      if (parentTeamId) {
        const parentTeam = await Team.findOne({
          _id: parentTeamId,
          tenantId: new mongoose.Types.ObjectId(tenantId),
        });

        if (!parentTeam) {
          throw new BadRequestError("Invalid parent team");
        }
      }

      // Create team
      const team = Team.build({
        name,
        description,
        tenantId: new mongoose.Types.ObjectId(tenantId),
        parentTeamId: parentTeamId
          ? new mongoose.Types.ObjectId(parentTeamId)
          : undefined,
        owner: new mongoose.Types.ObjectId(ownerId),
        settings,
      });

      // Add owner as first member with admin role
      team.members = [
        {
          userId: new mongoose.Types.ObjectId(ownerId),
          role: "admin",
          joinedAt: new Date(),
        },
      ];

      await team.save({ session });
      await session.commitTransaction();
      return team;
    } catch (error) {
      await session.abortTransaction();
      throw error;
    } finally {
      session.endSession();
    }
  }

  // Get team by ID
  static async getTeamById(teamId: string, tenantId: string): Promise<TeamDoc> {
    const team = await Team.findOne({
      _id: teamId,
      tenantId: new mongoose.Types.ObjectId(tenantId),
    })
      .populate("owner", "firstName lastName email")
      .populate("members.userId", "firstName lastName email");

    if (!team) {
      throw new NotFoundError("Team not found");
    }

    return team;
  }

  // Get all teams for a tenant
  static async getTeamsByTenant(
    tenantId: string,
    parentTeamId?: string
  ): Promise<TeamDoc[]> {
    const query: any = { tenantId: new mongoose.Types.ObjectId(tenantId) };
    if (parentTeamId !== undefined) {
      query.parentTeamId = parentTeamId
        ? new mongoose.Types.ObjectId(parentTeamId)
        : null;
    }

    return Team.find(query)
      .populate("owner", "firstName lastName email")
      .populate("members.userId", "firstName lastName email");
  }

  // Update team
  static async updateTeam(
    teamId: string,
    tenantId: string,
    updates: {
      name?: string;
      description?: string;
      settings?: {
        resourceQuota?: {
          maxMembers?: number;
          maxStorage?: number;
          maxProjects?: number;
        };
        visibility?: string;
        allowResourceSharing?: boolean;
      };
    }
  ): Promise<TeamDoc> {
    const team = await Team.findOneAndUpdate(
      {
        _id: teamId,
        tenantId: new mongoose.Types.ObjectId(tenantId),
      },
      { $set: updates },
      { new: true, runValidators: true }
    );

    if (!team) {
      throw new NotFoundError("Team not found");
    }

    return team;
  }

  // Add member to team
  static async addTeamMember(
    teamId: string,
    tenantId: string,
    userId: string,
    role: string
  ): Promise<TeamDoc> {
    const session = await mongoose.startSession();
    session.startTransaction();

    try {
      // Verify user exists and belongs to tenant
      const user = await User.findOne({
        _id: userId,
        tenantId: new mongoose.Types.ObjectId(tenantId),
      });

      if (!user) {
        throw new BadRequestError("Invalid user");
      }

      const team = await Team.findOne({
        _id: teamId,
        tenantId: new mongoose.Types.ObjectId(tenantId),
      });

      if (!team) {
        throw new NotFoundError("Team not found");
      }

      // Check if user is already a member
      if (team.members.some((member) => member.userId.toString() === userId)) {
        throw new BadRequestError("User is already a team member");
      }

      // Check resource quota
      if (team.members.length >= team.settings.resourceQuota!.maxMembers) {
        throw new BadRequestError("Team has reached maximum member limit");
      }

      team.members.push({
        userId: new mongoose.Types.ObjectId(userId),
        role,
        joinedAt: new Date(),
      });

      await team.save({ session });
      await session.commitTransaction();
      return team;
    } catch (error) {
      await session.abortTransaction();
      throw error;
    } finally {
      session.endSession();
    }
  }

  // Remove member from team
  static async removeTeamMember(
    teamId: string,
    tenantId: string,
    userId: string
  ): Promise<TeamDoc> {
    const team = await Team.findOne({
      _id: teamId,
      tenantId: new mongoose.Types.ObjectId(tenantId),
    });

    if (!team) {
      throw new NotFoundError("Team not found");
    }

    // Cannot remove owner
    if (team.owner.toString() === userId) {
      throw new ForbiddenError("Cannot remove team owner");
    }

    const updatedTeam = await Team.findOneAndUpdate(
      {
        _id: teamId,
        tenantId: new mongoose.Types.ObjectId(tenantId),
      },
      { $pull: { members: { userId: new mongoose.Types.ObjectId(userId) } } },
      { new: true, runValidators: true }
    );

    if (!updatedTeam) {
      throw new NotFoundError("Team not found");
    }

    return updatedTeam;
  }

  // Update member role
  static async updateMemberRole(
    teamId: string,
    tenantId: string,
    userId: string,
    newRole: string
  ): Promise<TeamDoc> {
    const team = await Team.findOne({
      _id: teamId,
      tenantId: new mongoose.Types.ObjectId(tenantId),
      "members.userId": new mongoose.Types.ObjectId(userId),
    });

    if (!team) {
      throw new NotFoundError("Team or member not found");
    }

    // Cannot change owner's role
    if (team.owner.toString() === userId) {
      throw new ForbiddenError("Cannot change team owner's role");
    }

    const updatedTeam = await Team.findOneAndUpdate(
      {
        _id: teamId,
        tenantId: new mongoose.Types.ObjectId(tenantId),
        "members.userId": new mongoose.Types.ObjectId(userId),
      },
      { $set: { "members.$.role": newRole } },
      { new: true, runValidators: true }
    );

    if (!updatedTeam) {
      throw new NotFoundError("Team or member not found");
    }

    return updatedTeam;
  }

  // Delete team
  static async deleteTeam(teamId: string, tenantId: string): Promise<void> {
    const session = await mongoose.startSession();
    session.startTransaction();

    try {
      // Check for child teams
      const childTeams = await Team.find({
        parentTeamId: new mongoose.Types.ObjectId(teamId),
      });

      if (childTeams.length > 0) {
        throw new BadRequestError("Cannot delete team with child teams");
      }

      const result = await Team.deleteOne({
        _id: teamId,
        tenantId: new mongoose.Types.ObjectId(tenantId),
      }).session(session);

      if (result.deletedCount === 0) {
        throw new NotFoundError("Team not found");
      }

      await session.commitTransaction();
    } catch (error) {
      await session.abortTransaction();
      throw error;
    } finally {
      session.endSession();
    }
  }
  // Get team hierarchy
  static async getTeamHierarchy(
    tenantId: string,
    rootTeamId?: string
  ): Promise<TeamHierarchy | TeamHierarchy[]> {
    const query: any = { tenantId: new mongoose.Types.ObjectId(tenantId) };
    if (rootTeamId) {
      query._id = new mongoose.Types.ObjectId(rootTeamId);
    } else {
      query.parentTeamId = null; // Get root level teams
    }

    const teams = await Team.find(query)
      .populate("owner", "firstName lastName email")
      .populate("members.userId", "firstName lastName email");

    const buildHierarchy = async (
      parentTeam: TeamDoc
    ): Promise<TeamHierarchy> => {
      const children = await Team.find({
        tenantId: new mongoose.Types.ObjectId(tenantId),
        parentTeamId: parentTeam._id,
      })
        .populate("owner", "firstName lastName email")
        .populate("members.userId", "firstName lastName email");

      const childHierarchies: TeamHierarchy[] = await Promise.all(
        children.map((child) => buildHierarchy(child))
      );

      return {
        ...parentTeam.toJSON(),
        children: childHierarchies,
      };
    };

    if (rootTeamId) {
      if (teams.length === 0) {
        throw new NotFoundError("Team not found");
      }
      return buildHierarchy(teams[0]);
    }

    return Promise.all(teams.map((team) => buildHierarchy(team)));
  }
}
