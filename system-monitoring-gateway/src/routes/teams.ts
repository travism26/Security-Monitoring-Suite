import express, { Request, Response } from "express";
import { TeamService } from "../services/team.service";
import { requireAuth } from "../middleware/require-auth";
import { validateRequest } from "../middleware/validate-request";
import { body } from "express-validator";

const router = express.Router();

// Create team validation
const createTeamValidation = [
  body("name").trim().notEmpty().withMessage("Team name is required"),
  body("description").optional().trim(),
  body("parentTeamId")
    .optional()
    .isMongoId()
    .withMessage("Invalid parent team ID"),
  body("settings")
    .optional()
    .isObject()
    .withMessage("Settings must be an object"),
  body("settings.resourceQuota")
    .optional()
    .isObject()
    .withMessage("Resource quota must be an object"),
  body("settings.visibility")
    .optional()
    .isIn(["private", "public"])
    .withMessage("Visibility must be either private or public"),
  body("settings.allowResourceSharing")
    .optional()
    .isBoolean()
    .withMessage("Allow resource sharing must be a boolean"),
];

// Create team
router.post(
  "/gateway/api/teams",
  requireAuth,
  createTeamValidation,
  validateRequest,
  async (req: Request, res: Response) => {
    const { name, description, parentTeamId, settings } = req.body;

    const team = await TeamService.createTeam(
      name,
      req.currentUser!.tenantId,
      req.currentUser!.id,
      description,
      parentTeamId,
      settings
    );

    res.status(201).send(team);
  }
);

// Get team by ID
router.get(
  "/gateway/api/teams/:teamId",
  requireAuth,
  async (req: Request, res: Response) => {
    const team = await TeamService.getTeamById(
      req.params.teamId,
      req.currentUser!.tenantId
    );
    res.send(team);
  }
);

// Get all teams for tenant
router.get(
  "/gateway/api/teams",
  requireAuth,
  async (req: Request, res: Response) => {
    const { parentTeamId } = req.query;
    const teams = await TeamService.getTeamsByTenant(
      req.currentUser!.tenantId,
      parentTeamId as string | undefined
    );
    res.send(teams);
  }
);

// Update team validation
const updateTeamValidation = [
  body("name")
    .optional()
    .trim()
    .notEmpty()
    .withMessage("Team name cannot be empty"),
  body("description").optional().trim(),
  body("settings")
    .optional()
    .isObject()
    .withMessage("Settings must be an object"),
  body("settings.resourceQuota")
    .optional()
    .isObject()
    .withMessage("Resource quota must be an object"),
  body("settings.visibility")
    .optional()
    .isIn(["private", "public"])
    .withMessage("Visibility must be either private or public"),
  body("settings.allowResourceSharing")
    .optional()
    .isBoolean()
    .withMessage("Allow resource sharing must be a boolean"),
];

// Update team
router.patch(
  "/gateway/api/teams/:teamId",
  requireAuth,
  updateTeamValidation,
  validateRequest,
  async (req: Request, res: Response) => {
    const team = await TeamService.updateTeam(
      req.params.teamId,
      req.currentUser!.tenantId,
      req.body
    );
    res.send(team);
  }
);

// Add member validation
const addMemberValidation = [
  body("userId").isMongoId().withMessage("Valid user ID is required"),
  body("role")
    .isIn(["admin", "member"])
    .withMessage("Role must be either admin or member"),
];

// Add team member
router.post(
  "/gateway/api/teams/:teamId/members",
  requireAuth,
  addMemberValidation,
  validateRequest,
  async (req: Request, res: Response) => {
    const { userId, role } = req.body;
    const team = await TeamService.addTeamMember(
      req.params.teamId,
      req.currentUser!.tenantId,
      userId,
      role
    );
    res.send(team);
  }
);

// Remove team member
router.delete(
  "/gateway/api/teams/:teamId/members/:userId",
  requireAuth,
  async (req: Request, res: Response) => {
    const team = await TeamService.removeTeamMember(
      req.params.teamId,
      req.currentUser!.tenantId,
      req.params.userId
    );
    res.send(team);
  }
);

// Update member role validation
const updateMemberRoleValidation = [
  body("role")
    .isIn(["admin", "member"])
    .withMessage("Role must be either admin or member"),
];

// Update member role
router.patch(
  "/gateway/api/teams/:teamId/members/:userId/role",
  requireAuth,
  updateMemberRoleValidation,
  validateRequest,
  async (req: Request, res: Response) => {
    const team = await TeamService.updateMemberRole(
      req.params.teamId,
      req.currentUser!.tenantId,
      req.params.userId,
      req.body.role
    );
    res.send(team);
  }
);

// Delete team
router.delete(
  "/gateway/api/teams/:teamId",
  requireAuth,
  async (req: Request, res: Response) => {
    await TeamService.deleteTeam(req.params.teamId, req.currentUser!.tenantId);
    res.status(204).send();
  }
);

// Get team hierarchy
router.get(
  "/gateway/api/teams/hierarchy",
  requireAuth,
  async (req: Request, res: Response) => {
    const { rootTeamId } = req.query;
    const hierarchy = await TeamService.getTeamHierarchy(
      req.currentUser!.tenantId,
      rootTeamId as string | undefined
    );
    res.send(hierarchy);
  }
);

export { router as teamsRouter };
