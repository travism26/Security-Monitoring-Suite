// Mocked data
const mockUsers = [{ id: "1", name: "Test User", email: "test@example.com" }];

const mockTeams = [
  { id: "1", name: "Team A" },
  { id: "2", name: "Team B" },
  { id: "3", name: "Team C" },
];

const mockApiKeys = [
  {
    id: "1",
    name: "Production Server",
    key: "prod_abcdefghijklmnop",
    createdAt: "2023-06-15T10:00:00Z",
  },
  {
    id: "2",
    name: "Development Server",
    key: "dev_qrstuvwxyz123456",
    createdAt: "2023-06-16T11:30:00Z",
  },
];

const mockEvents = [
  {
    id: 1,
    timestamp: "2023-06-10 10:30:15",
    type: "Login",
    source: "192.168.1.100",
    description: "Successful login",
  },
  {
    id: 2,
    timestamp: "2023-06-10 10:35:22",
    type: "File Access",
    source: "192.168.1.101",
    description: "Unauthorized file access attempt",
  },
  {
    id: 3,
    timestamp: "2023-06-10 10:40:05",
    type: "Network",
    source: "192.168.1.102",
    description: "Unusual outbound traffic detected",
  },
];

const mockThreats = [
  { name: "Malware", count: 15, severity: 70 },
  { name: "Phishing", count: 8, severity: 60 },
  { name: "DDoS", count: 3, severity: 40 },
];

const mockSystemHealth = [
  { name: "Firewall", status: "Operational" },
  { name: "IDS", status: "Operational" },
  { name: "Log Server", status: "Down" },
  { name: "Email Filter", status: "Operational" },
];

const mockAlerts = [
  {
    id: 1,
    message: "High CPU usage detected",
    timestamp: "2023-06-17T14:30:00Z",
  },
  {
    id: 2,
    message: "Unusual network activity",
    timestamp: "2023-06-17T15:45:00Z",
  },
  {
    id: 3,
    message: "Failed login attempts",
    timestamp: "2023-06-17T16:20:00Z",
  },
];

// API Service
export const apiService = {
  // Auth
  login: async (email: string, password: string) => {
    const user = mockUsers.find((u) => u.email === email);
    if (user && password === "password123") {
      return user;
    }
    throw new Error("Invalid credentials");
  },

  signup: async (email: string, password: string) => {
    const newUser = {
      id: String(mockUsers.length + 1),
      name: "New User",
      email,
    };
    mockUsers.push(newUser);
    return newUser;
  },

  // Teams
  getTeams: async () => {
    return mockTeams;
  },

  switchTeam: async (teamId: string) => {
    const team = mockTeams.find((t) => t.id === teamId);
    if (team) {
      return team;
    }
    throw new Error("Team not found");
  },

  // API Keys
  getApiKeys: async () => {
    return mockApiKeys;
  },

  createApiKey: async (name: string) => {
    const newKey = {
      id: String(mockApiKeys.length + 1),
      name,
      key: `new_${Math.random().toString(36).substr(2, 16)}`,
      createdAt: new Date().toISOString(),
    };
    mockApiKeys.push(newKey);
    return newKey;
  },

  deleteApiKey: async (id: string) => {
    const index = mockApiKeys.findIndex((key) => key.id === id);
    if (index !== -1) {
      mockApiKeys.splice(index, 1);
      return true;
    }
    throw new Error("API key not found");
  },

  // Dashboard Data
  getEvents: async () => {
    return mockEvents;
  },

  getThreats: async () => {
    return mockThreats;
  },

  getSystemHealth: async () => {
    return mockSystemHealth;
  },

  getAlerts: async () => {
    return mockAlerts;
  },
};
