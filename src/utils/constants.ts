// utils/constants.ts
import type { ProfessorKey } from '../interfaces/Professor';

// List of possible options for the SubField filter
export const optionsList = [
  "Artificial Intelligence",
  "Software Engineering",
  "Computer Security",
  "Databases",
  "Cryptography",
  "Programming Languages",
];

// Define custom column widths for each column
export const columnWidths: { [key in ProfessorKey]: string } = {
  FullName: '15%',
  University: '20%',
  JoinYear: '5%',
  SubField: '20%',
  Bachelors: '20%',
  Doctorate: '20%',
};