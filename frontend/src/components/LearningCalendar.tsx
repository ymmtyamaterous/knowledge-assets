"use client";

import type { CalendarDay } from "@/types/progress";

type Props = {
  days: CalendarDay[];
  year: number;
};

function getCellColor(count: number): string {
  if (count === 0) return "bg-slate-100";
  if (count === 1) return "bg-pink-100";
  if (count <= 3) return "bg-pink-300";
  return "bg-pink-500";
}

export default function LearningCalendar({ days, year }: Props) {
  // 日付 → 完了数のマップ
  const countMap: Record<string, number> = {};
  for (const d of days) {
    countMap[d.date] = d.count;
  }

  // 1月1日から52週（364日）+ 余りをグリッドとして構築
  const startDate = new Date(year, 0, 1); // 1/1
  // 週の始まりを日曜にそろえる
  const startDow = startDate.getDay(); // 0=Sun
  const gridStart = new Date(startDate);
  gridStart.setDate(gridStart.getDate() - startDow);

  // 53週分 × 7日 = 371日分のセルを生成
  const totalWeeks = 53;
  const cells: { date: string; count: number; inYear: boolean }[][] = [];
  for (let w = 0; w < totalWeeks; w++) {
    const week: { date: string; count: number; inYear: boolean }[] = [];
    for (let d = 0; d < 7; d++) {
      const cur = new Date(gridStart);
      cur.setDate(cur.getDate() + w * 7 + d);
      const key = cur.toISOString().slice(0, 10);
      week.push({
        date: key,
        count: countMap[key] ?? 0,
        inYear: cur.getFullYear() === year,
      });
    }
    cells.push(week);
  }

  // 月ラベル（各月の最初の週を検出）
  const monthLabels: { week: number; label: string }[] = [];
  const seenMonths = new Set<number>();
  for (let w = 0; w < totalWeeks; w++) {
    const firstDay = cells[w][0];
    if (!firstDay.inYear) continue;
    const month = new Date(firstDay.date).getMonth();
    if (!seenMonths.has(month)) {
      seenMonths.add(month);
      monthLabels.push({
        week: w,
        label: `${month + 1}月`,
      });
    }
  }

  const CELL_SIZE = 12; // px
  const GAP = 2;
  const DOW_LABELS = ["", "月", "", "水", "", "金", ""];

  return (
    <div className="overflow-x-auto">
      <div className="inline-block">
        {/* 月ラベル */}
        <div className="mb-1 flex" style={{ paddingLeft: 20 }}>
          {Array.from({ length: totalWeeks }, (_, w) => {
            const label = monthLabels.find((m) => m.week === w);
            return (
              <div
                key={w}
                style={{ width: CELL_SIZE + GAP, flexShrink: 0 }}
                className="text-center text-[10px] text-slate-400"
              >
                {label ? label.label : ""}
              </div>
            );
          })}
        </div>

        {/* グリッド本体 */}
        <div className="flex gap-[2px]">
          {/* 曜日ラベル */}
          <div
            className="flex flex-col justify-between"
            style={{ width: 16, gap: GAP, marginRight: GAP }}
          >
            {DOW_LABELS.map((label, i) => (
              <div
                key={i}
                style={{ height: CELL_SIZE }}
                className="flex items-center justify-end text-[10px] text-slate-400"
              >
                {label}
              </div>
            ))}
          </div>

          {/* 週ごとのカラム */}
          {cells.map((week, w) => (
            <div key={w} className="flex flex-col gap-[2px]">
              {week.map((cell, d) => (
                <div
                  key={d}
                  title={cell.inYear ? `${cell.date}: ${cell.count}件` : ""}
                  className={`rounded-sm transition-opacity ${cell.inYear ? getCellColor(cell.count) : "bg-transparent"}`}
                  style={{ width: CELL_SIZE, height: CELL_SIZE }}
                />
              ))}
            </div>
          ))}
        </div>

        {/* 凡例 */}
        <div className="mt-2 flex items-center gap-2 text-[10px] text-slate-400">
          <span>少ない</span>
          {["bg-slate-100", "bg-pink-100", "bg-pink-300", "bg-pink-500"].map((c) => (
            <div key={c} className={`${c} rounded-sm`} style={{ width: CELL_SIZE, height: CELL_SIZE }} />
          ))}
          <span>多い</span>
        </div>
      </div>
    </div>
  );
}
