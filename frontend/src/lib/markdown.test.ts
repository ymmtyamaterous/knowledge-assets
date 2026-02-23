import { convertMarkdownToHtml } from "@/lib/markdown";

describe("convertMarkdownToHtml", () => {
  it("renders markdown table with table tags", () => {
    const input = "| A | B |\n|---|---|\n| 1 | 2 |";
    const html = convertMarkdownToHtml(input);

    expect(html).toContain("<table");
    expect(html).toContain("<thead>");
    expect(html).toContain("<tbody>");
    expect(html).not.toContain("|---|---|");
  });

  it("renders ordered list", () => {
    const input = "1. first\n2. second\n3. third";
    const html = convertMarkdownToHtml(input);

    expect(html).toContain("<ol");
    expect(html).toContain("<li>first</li>");
    expect(html).toContain("<li>second</li>");
    expect(html).toContain("<li>third</li>");
  });

  it("renders unordered list", () => {
    const input = "- alpha\n- beta";
    const html = convertMarkdownToHtml(input);

    expect(html).toContain("<ul");
    expect(html).toContain("<li>alpha</li>");
    expect(html).toContain("<li>beta</li>");
  });
});
