// interfaces/Professor.ts

// Define the structure of the data for each professor
export interface Professor {
  FullName: string;
  University: string;
  JoinYear: string;
  SubField: string[];
  Bachelors: string;
  Doctorate: string;
}

// List of valid keys and corresponding TypeScript type
export const validKeys = ["FullName", "University", "JoinYear", "SubField", "Bachelors", "Doctorate"] as const;
export type ProfessorKey = typeof validKeys[number];