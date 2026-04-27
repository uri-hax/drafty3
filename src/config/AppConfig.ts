export const datasetFiles: Record<string, { csv: string; yaml: string; }> = {
  csprofs: {
    csv: "suggestions.csv",
    yaml: "csprofessors.yaml",
  },
  // students: {
  //   csv: "students.csv",
  //   yaml: "students.yaml",
  // },
};

export const datasetLabels: Record<string, string> = {
  csprofs: "CS Professors",
  //students: "Students",
};

export const editFiles: Record<string, string> = {
  csprofs: "edit-history.csv",
  //students: "students-edit-history.csv",
};