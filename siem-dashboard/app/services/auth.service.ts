export interface User {
  id: string;
  email: string;
  firstName: string;
  lastName: string;
  role: string;
}

interface AuthResponse {
  user: User;
  token: string;
}

interface SignupData {
  email: string;
  password: string;
  firstName: string;
  lastName: string;
  tenantId: string;
}

export class AuthService {
  // UPDATE_BACKEND_BASE_URL: Replace with your k8s nodeport base path
  private static gatewayUrl = "/gateway/api/v1";

  private static async logApiCall(
    endpoint: string,
    status: number,
    error?: any
  ) {
    console.log(`[Auth Service] API Call to ${endpoint}`, {
      timestamp: new Date().toISOString(),
      status,
      error: error ? JSON.stringify(error) : undefined,
    });
  }

  static async login(email: string, password: string): Promise<AuthResponse> {
    const endpoint = `${this.gatewayUrl}/auth/login`;
    console.log(`[Auth Service] Attempting login for user ${email}`);

    try {
      const response: Response = await fetch(endpoint, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        credentials: "include",
        body: JSON.stringify({ email, password }),
      });

      if (!response.ok) {
        const error = await response.json();
        this.logApiCall(endpoint, response.status, error);
        throw new Error(error.message || "Failed to login");
      }

      this.logApiCall(endpoint, response.status);
      const data: AuthResponse = await response.json();
      console.log("[Auth Service] Login successful, redirecting to dashboard");
      return data;
    } catch (error) {
      console.error("[Auth Service] Login error:", error);
      throw error;
    }
  }

  static async signup(
    email: string,
    password: string,
    firstName: string,
    lastName: string
  ): Promise<AuthResponse> {
    const endpoint = `${this.gatewayUrl}/auth/register`;
    console.log(`[Auth Service] Attempting signup for user ${email}`);

    // Use a default tenant ID for now - in a real app this would be handled differently
    const signupData = {
      email,
      password,
      firstName,
      lastName,
      tenantId: "default-tenant",
    };

    try {
      const response: Response = await fetch(endpoint, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        credentials: "include",
        body: JSON.stringify(signupData),
      });

      if (!response.ok) {
        const error = await response.json();
        this.logApiCall(endpoint, response.status, error);
        throw new Error(error.message || "Failed to signup");
      }

      this.logApiCall(endpoint, response.status);
      const data: AuthResponse = await response.json();
      console.log("[Auth Service] Signup successful, redirecting to dashboard");
      return data;
    } catch (error) {
      console.error("[Auth Service] Signup error:", error);
      throw error;
    }
  }

  static async getCurrentUser(): Promise<User | null> {
    const endpoint = `${this.gatewayUrl}/users/me`;
    console.log("[Auth Service] Fetching current user data");

    try {
      const response = await fetch(endpoint, {
        method: "GET",
        credentials: "include",
        headers: {
          Accept: "application/json",
          "Content-Type": "application/json",
        },
      });

      if (!response.ok) {
        this.logApiCall(endpoint, response.status);
        console.log("[Auth Service] Failed to fetch current user");
        return null;
      }

      this.logApiCall(endpoint, response.status);
      const userData = await response.json();
      console.log(
        "[Auth Service] Successfully fetched current user data:",
        userData
      );
      return userData;
    } catch (error) {
      console.error("[Auth Service] Error fetching current user:", error);
      return null;
    }
  }
}
