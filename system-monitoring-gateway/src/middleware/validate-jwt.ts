import { Request, Response, NextFunction } from "express";
import { NotAuthorizedError } from "../errors";
import { JWTService } from "../services/jwt.service";
import { UserPayload } from "../types/auth";

declare global {
  namespace Express {
    interface Request {
      currentUser?: UserPayload;
    }
  }
}

export const validateJWT = (
  req: Request,
  res: Response,
  next: NextFunction
) => {
  const authHeader = req.headers.authorization;
  console.log("authHeader", authHeader);
  console.log(`req.currentUser`, req?.currentUser);
  console.log(`req.session`, req?.session);

  // if (!authHeader || !authHeader.startsWith("Bearer ")) {
  //   console.log("No auth header");
  //   throw new NotAuthorizedError("Authorization header missing or invalid");
  // }

  if (!req.session?.jwt) {
    // return next();
    throw new NotAuthorizedError("Not Authorized");
  }

  try {
    const token = req.session.jwt;
    const payload = JWTService.verifyToken(token);
    console.log("payload", payload);
    req.currentUser = payload;
    next();
  } catch (error) {
    throw new NotAuthorizedError("Invalid token");
  }
};
