export type GlossaryTerm = {
  id: string;
  term: string;
  reading: string;
  definition: string;
};

export type GlossaryResponse = {
  terms: GlossaryTerm[];
};
