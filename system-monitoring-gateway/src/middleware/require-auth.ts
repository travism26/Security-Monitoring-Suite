import { Request, Response, NextFunction } from "express";
import { NotAuthorizedError } from "../errors";

declare global {
  namespace Express {
    interface Request {
      currentUser?: {
        id: string;
        email: string;
        tenantId: string;
        role: string;
      };
    }
  }
}

export const requireAuth = (
  req: Request,
  res: Response,
  next: NextFunction
) => {
  if (!req.currentUser) {
    throw new NotAuthorizedError("Not authorized");
  }

  next();
};
