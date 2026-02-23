import { fetchCourses } from "@/lib/api";

describe("fetchCourses", () => {
  const originalFetch = global.fetch;

  afterEach(() => {
    global.fetch = originalFetch;
    jest.restoreAllMocks();
  });

  it("returns empty array when courses is undefined", async () => {
    global.fetch = jest.fn().mockResolvedValue({
      ok: true,
      json: async () => ({}),
    } as Response);

    await expect(fetchCourses()).resolves.toEqual([]);
  });

  it("throws on non-ok response", async () => {
    global.fetch = jest.fn().mockResolvedValue({ ok: false } as Response);

    await expect(fetchCourses()).rejects.toThrow("failed to fetch courses");
  });
});
