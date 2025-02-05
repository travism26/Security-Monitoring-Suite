import jwt from "jsonwebtoken";

const JWT_SECRET = process.env.JWT_SECRET || "your-secret-key"; // In production, always use environment variable

interface UserPayload {
  id: string;
  email: string;
  tenantId?: string;
  role: string;
}

export class JWTService {
  static generateToken(payload: UserPayload): string {
    return jwt.sign(payload, JWT_SECRET, {
      expiresIn: "24h", // Token expires in 24 hours
    });
  }

  static verifyToken(token: string): UserPayload {
    try {
      const payload = jwt.verify(token, JWT_SECRET) as UserPayload;
      return payload;
    } catch (error) {
      throw new Error("Invalid token");
    }
  }
}
