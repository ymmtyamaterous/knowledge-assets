export type GlossaryTag = {
  id: string;
  name: string;
  createdAt: string;
};

export type GlossaryTerm = {
  id: string;
  term: string;
  reading: string;
  definition: string;
  tags: GlossaryTag[];
};

export type GlossaryResponse = {
  terms: GlossaryTerm[];
};

export type GlossaryTagsResponse = {
  tags: GlossaryTag[];
};
