// UserPayload represents the authenticated user information
// tenantId is optional when multi-tenancy is disabled
export interface UserPayload {
  id: string;
  email: string;
  tenantId?: string;
  role: string;
}
