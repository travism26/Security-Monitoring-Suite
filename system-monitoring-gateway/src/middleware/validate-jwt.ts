import { Request, Response, NextFunction } from "express";
import { NotAuthorizedError } from "../errors";
import { JWTService } from "../services/jwt.service";

export const validateJWT = (
  req: Request,
  res: Response,
  next: NextFunction
) => {
  const authHeader = req.headers.authorization;

  if (!authHeader || !authHeader.startsWith("Bearer ")) {
    throw new NotAuthorizedError("Authorization header missing or invalid");
  }

  try {
    const token = authHeader.split(" ")[1];
    const payload = JWTService.verifyToken(token);
    req.currentUser = payload;
    next();
  } catch (error) {
    throw new NotAuthorizedError("Invalid token");
  }
};
