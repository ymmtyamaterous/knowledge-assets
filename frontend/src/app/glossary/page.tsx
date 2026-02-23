import { fetchGlossary } from "@/lib/api";

export default async function GlossaryPage() {
  const terms = await fetchGlossary().catch(() => []);

  return (
    <main>
      <h1 className="mb-6 text-2xl font-bold text-slate-800">用語辞典</h1>

      {terms.length === 0 ? (
        <p className="text-sm text-slate-500">用語データを取得できませんでした。</p>
      ) : (
        <div className="space-y-3">
          {terms.map((term) => (
            <div
              key={term.id}
              className="rounded-xl border border-slate-200 bg-white p-4 shadow-sm"
            >
              <div className="flex items-baseline gap-2">
                <h2 className="text-lg font-semibold text-slate-800">{term.term}</h2>
                {term.reading && (
                  <span className="text-sm text-slate-400">（{term.reading}）</span>
                )}
              </div>
              <p className="mt-2 text-sm text-slate-600">{term.definition}</p>
            </div>
          ))}
        </div>
      )}
    </main>
  );
}
