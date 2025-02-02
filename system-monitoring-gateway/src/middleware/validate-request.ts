import { Request, Response, NextFunction } from "express";
import { validationResult } from "express-validator";
import { BadRequestError } from "../errors";

export const validateRequest = (
  req: Request,
  res: Response,
  next: NextFunction
) => {
  const errors = validationResult(req);

  if (!errors.isEmpty()) {
    const formattedErrors = errors.array().map((error) => ({
      message: error.msg,
      field: error.type === "field" ? error.path : undefined,
    }));
    throw new BadRequestError(
      `Validation failed: ${formattedErrors
        .map((e) => `${e.field}: ${e.message}`)
        .join(", ")}`
    );
  }

  next();
};
