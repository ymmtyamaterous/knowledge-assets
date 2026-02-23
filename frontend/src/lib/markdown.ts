function escapeHtml(value: string): string {
  return value
    .replaceAll("&", "&amp;")
    .replaceAll("<", "&lt;")
    .replaceAll(">", "&gt;")
    .replaceAll('"', "&quot;")
    .replaceAll("'", "&#39;");
}

function isTableSeparator(line: string): boolean {
  return /^\|?\s*:?-{3,}:?\s*(\|\s*:?-{3,}:?\s*)+\|?$/.test(line.trim());
}

function renderInlineMarkdown(text: string): string {
  return escapeHtml(text)
    .replace(/\*\*(.+?)\*\*/g, "<strong>$1</strong>")
    .replace(/\*(.+?)\*/g, "<em>$1</em>");
}

export function convertMarkdownToHtml(content: string): string {
  const lines = content.split("\n");
  const html: string[] = [];

  let i = 0;
  while (i < lines.length) {
    const line = lines[i];

    if (/^```/.test(line)) {
      const codeLines: string[] = [];
      i += 1;
      while (i < lines.length && !/^```/.test(lines[i])) {
        codeLines.push(lines[i]);
        i += 1;
      }
      html.push(`<pre class='bg-slate-100 rounded p-3 text-sm overflow-x-auto my-3'><code>${escapeHtml(codeLines.join("\n"))}</code></pre>`);
      i += 1;
      continue;
    }

    if (/^\|.*\|$/.test(line) && i + 1 < lines.length && isTableSeparator(lines[i + 1])) {
      const headerCells = line
        .split("|")
        .filter((cell) => cell.trim() !== "")
        .map((cell) => `<th class='border border-slate-300 bg-slate-50 px-3 py-1 text-left text-sm font-semibold'>${renderInlineMarkdown(cell.trim())}</th>`)
        .join("");

      const rows: string[] = [];
      i += 2;
      while (i < lines.length && /^\|.*\|$/.test(lines[i])) {
        const rowCells = lines[i]
          .split("|")
          .filter((cell) => cell.trim() !== "")
          .map((cell) => `<td class='border border-slate-300 px-3 py-1 text-sm'>${renderInlineMarkdown(cell.trim())}</td>`)
          .join("");
        rows.push(`<tr>${rowCells}</tr>`);
        i += 1;
      }

      html.push(`<div class='my-3 overflow-x-auto'><table class='w-full border-collapse'><thead><tr>${headerCells}</tr></thead><tbody>${rows.join("")}</tbody></table></div>`);
      continue;
    }

    if (/^\d+\.\s+/.test(line)) {
      const items: string[] = [];
      while (i < lines.length && /^\d+\.\s+/.test(lines[i])) {
        items.push(lines[i].replace(/^\d+\.\s+/, "").trim());
        i += 1;
      }
      html.push(`<ol class='ml-6 list-decimal space-y-1'>${items.map((item) => `<li>${renderInlineMarkdown(item)}</li>`).join("")}</ol>`);
      continue;
    }

    if (/^-\s+/.test(line)) {
      const items: string[] = [];
      while (i < lines.length && /^-\s+/.test(lines[i])) {
        items.push(lines[i].replace(/^-\s+/, "").trim());
        i += 1;
      }
      html.push(`<ul class='ml-6 list-disc space-y-1'>${items.map((item) => `<li>${renderInlineMarkdown(item)}</li>`).join("")}</ul>`);
      continue;
    }

    if (/^###\s+/.test(line)) {
      html.push(`<h3 class='text-lg font-semibold mt-4 mb-1'>${renderInlineMarkdown(line.replace(/^###\s+/, ""))}</h3>`);
      i += 1;
      continue;
    }
    if (/^##\s+/.test(line)) {
      html.push(`<h2 class='text-xl font-bold mt-6 mb-2 text-pink-500'>${renderInlineMarkdown(line.replace(/^##\s+/, ""))}</h2>`);
      i += 1;
      continue;
    }
    if (/^#\s+/.test(line)) {
      html.push(`<h1 class='text-2xl font-bold mt-6 mb-2'>${renderInlineMarkdown(line.replace(/^#\s+/, ""))}</h1>`);
      i += 1;
      continue;
    }

    if (line.trim() === "") {
      html.push("");
      i += 1;
      continue;
    }

    html.push(`<p class='mt-3'>${renderInlineMarkdown(line)}</p>`);
    i += 1;
  }

  return html.join("\n");
}
