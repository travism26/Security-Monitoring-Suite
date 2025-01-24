import { Request, Response, NextFunction } from "express";
import jwt from "jsonwebtoken";
import { NotAuthorizedError } from "../errors/not-authorized-error";

interface UserPayload {
  id: string;
  tenantId: string;
  role: string;
}

// Augment the Express Request type to include currentUser
declare global {
  namespace Express {
    interface Request {
      currentUser?: UserPayload;
    }
  }
}

export const requireAuth = (
  req: Request,
  res: Response,
  next: NextFunction
) => {
  if (!req.currentUser) {
    throw new NotAuthorizedError();
  }

  next();
};

export const validateJWT = (
  req: Request,
  res: Response,
  next: NextFunction
) => {
  const authHeader = req.headers.authorization;

  if (!authHeader?.startsWith("Bearer ")) {
    throw new NotAuthorizedError();
  }

  try {
    const token = authHeader.split(" ")[1];
    const payload = jwt.verify(token, process.env.JWT_KEY!) as UserPayload;

    req.currentUser = payload;
    next();
  } catch (err) {
    throw new NotAuthorizedError();
  }
};
