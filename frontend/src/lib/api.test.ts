import { fetchCourses, loginUser, fetchGlossary } from "@/lib/api";

// localStorage のスタブ（token管理のため）
const localStorageMock = (() => {
  const store: Record<string, string> = {};
  return {
    getItem: (key: string) => store[key] ?? null,
    setItem: (key: string, value: string) => { store[key] = value; },
    removeItem: (key: string) => { delete store[key]; },
    clear: () => { Object.keys(store).forEach((k) => delete store[k]); },
  };
})();
Object.defineProperty(global, "localStorage", { value: localStorageMock });

const mockFetch = (ok: boolean, body: unknown) =>
  jest.fn().mockResolvedValue({
    ok,
    json: async () => body,
  } as Response);

describe("fetchCourses", () => {
  const originalFetch = global.fetch;

  afterEach(() => {
    global.fetch = originalFetch;
    jest.restoreAllMocks();
  });

  it("returns courses array on success", async () => {
    global.fetch = mockFetch(true, {
      courses: [{ id: "c1", title: "FP3", description: "", difficulty: "初級", estimatedHour: 10 }],
    });
    const courses = await fetchCourses();
    expect(courses).toHaveLength(1);
    expect(courses[0].id).toBe("c1");
  });

  it("returns empty array when courses key is missing", async () => {
    global.fetch = mockFetch(true, {});
    await expect(fetchCourses()).resolves.toEqual([]);
  });

  it("throws on non-ok response", async () => {
    global.fetch = mockFetch(false, { error: "not found" });
    await expect(fetchCourses()).rejects.toThrow("not found");
  });
});

describe("loginUser", () => {
  const originalFetch = global.fetch;

  afterEach(() => {
    global.fetch = originalFetch;
  });

  it("returns AuthResponse on success", async () => {
    global.fetch = mockFetch(true, { token: "tok123", user: { id: "u1", email: "a@a.com", username: "alice" } });
    const res = await loginUser("a@a.com", "password");
    expect(res.token).toBe("tok123");
    expect(res.user.email).toBe("a@a.com");
  });

  it("throws on invalid credentials", async () => {
    global.fetch = mockFetch(false, { error: "invalid credentials" });
    await expect(loginUser("x@x.com", "wrong")).rejects.toThrow("invalid credentials");
  });
});

describe("fetchGlossary", () => {
  const originalFetch = global.fetch;

  afterEach(() => {
    global.fetch = originalFetch;
  });

  it("returns terms array on success", async () => {
    global.fetch = mockFetch(true, {
      terms: [{ id: "g1", term: "複利", reading: "ふくり", definition: "元本に利息を加えた金額に再び利息がつく計算方式" }],
    });
    const terms = await fetchGlossary();
    expect(terms).toHaveLength(1);
    expect(terms[0].term).toBe("複利");
  });

  it("returns empty array when terms key is missing", async () => {
    global.fetch = mockFetch(true, {});
    await expect(fetchGlossary()).resolves.toEqual([]);
  });
});
