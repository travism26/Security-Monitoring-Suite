import express, { Request, Response } from "express";
import { requireAuth } from "../middleware/require-auth";
import { validateJWT } from "../middleware/validate-jwt";
import { JWTService } from "../services/jwt.service";
import { validateTenantConsistency } from "../middlewares/validate-tenant";
import { UserService } from "../services/user.service";

const router = express.Router();

// All routes require JWT authentication except registration and login
router.post("/gateway/register", async (req: Request, res: Response) => {
  const { email, password, firstName, lastName, tenantId } = req.body;

  try {
    // Check if user already exists
    const existingUser = await UserService.getUserByEmail(email);
    if (existingUser) {
      return res.status(400).send({ message: "Email already in use" });
    }

    const user = await UserService.createUser({
      email,
      password,
      firstName,
      lastName,
      tenantId,
    });

    // Generate token for immediate login after registration
    const token = JWTService.generateToken({
      id: user.id,
      email: user.email,
      tenantId: user.tenantId,
      role: user.role,
    });

    res.status(201).send({ user, token });
  } catch (error) {
    res.status(400).send({ message: "Invalid user data" });
  }
});

router.post("/gateway/login", async (req: Request, res: Response) => {
  const { email, password } = req.body;

  try {
    const user = await UserService.validateCredentials(email, password);
    if (!user) {
      return res.status(401).send({ message: "Invalid credentials" });
    }

    const token = JWTService.generateToken({
      id: user.id,
      email: user.email,
      tenantId: user.tenantId,
      role: user.role,
    });

    res.send({ user, token });
  } catch (error) {
    res.status(400).send({ message: "Invalid login attempt" });
  }
});

// Protected routes
router.use(validateJWT);
router.use(requireAuth);
router.use(validateTenantConsistency);

// Mount all protected routes under /api/v1/users
const usersRouter = express.Router();
usersRouter.use("/gateway/api/v1/users", router);

// Get current user profile
router.get("/me", async (req: Request, res: Response) => {
  const userId = req.currentUser!.id;
  const user = await UserService.getUserById(userId);

  if (!user) {
    return res.status(404).send({ message: "User not found" });
  }

  res.send(user);
});

// Update user profile
router.patch("/gateway/me", async (req: Request, res: Response) => {
  const userId = req.currentUser!.id;
  const { firstName, lastName } = req.body;

  const user = await UserService.updateUser(userId, { firstName, lastName });

  if (!user) {
    return res.status(404).send({ message: "User not found" });
  }

  res.send(user);
});

// Admin routes - require admin role
const requireAdmin = (req: Request, res: Response, next: Function) => {
  if (req.currentUser?.role !== "admin") {
    return res.status(403).send({ message: "Admin access required" });
  }
  next();
};

// List users (admin only)
router.get("/gateway/", requireAdmin, async (req: Request, res: Response) => {
  const tenantId = req.currentUser!.tenantId;
  const users = await UserService.listUsersByTenant(tenantId);
  res.send(users);
});

// Create user (admin only)
router.post("/gateway/", requireAdmin, async (req: Request, res: Response) => {
  const { email, password, firstName, lastName, role } = req.body;
  const tenantId = req.currentUser!.tenantId;

  try {
    const existingUser = await UserService.getUserByEmail(email);
    if (existingUser) {
      return res.status(400).send({ message: "Email already in use" });
    }

    const user = await UserService.createUser({
      email,
      password,
      firstName,
      lastName,
      role,
      tenantId,
    });

    res.status(201).send(user);
  } catch (error) {
    res.status(400).send({ message: "Invalid user data" });
  }
});

// Update user (admin only)
router.patch(
  "/gateway/:userId",
  requireAdmin,
  async (req: Request, res: Response) => {
    const { userId } = req.params;
    const { firstName, lastName, role, status } = req.body;
    const tenantId = req.currentUser!.tenantId;

    // Verify user belongs to tenant
    const user = await UserService.getUserById(userId);
    if (!user || user.tenantId !== tenantId) {
      return res.status(404).send({ message: "User not found" });
    }

    const updatedUser = await UserService.updateUser(userId, {
      firstName,
      lastName,
      role,
      status,
    });

    res.send(updatedUser);
  }
);

// Password reset request
router.post("/gateway/forgot-password", async (req: Request, res: Response) => {
  const { email } = req.body;

  const success = await UserService.initiatePasswordReset(email);
  // Always return success to prevent email enumeration
  res.send({ message: "If the email exists, a reset link will be sent" });
});

// Reset password with token
router.post("/gateway/reset-password", async (req: Request, res: Response) => {
  const { token, newPassword } = req.body;

  const success = await UserService.resetPassword(token, newPassword);
  if (!success) {
    return res.status(400).send({ message: "Invalid or expired reset token" });
  }

  res.send({ message: "Password successfully reset" });
});

// Verify email
router.get(
  "/gateway/verify-email/:token",
  async (req: Request, res: Response) => {
    const { token } = req.params;

    const success = await UserService.verifyEmail(token);
    if (!success) {
      return res.status(400).send({ message: "Invalid verification token" });
    }

    res.send({ message: "Email successfully verified" });
  }
);

export { usersRouter };
