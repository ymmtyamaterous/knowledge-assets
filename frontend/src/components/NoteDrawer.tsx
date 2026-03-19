"use client";

import { useEffect, useRef } from "react";

type Props = {
  isOpen: boolean;
  onClose: () => void;
  noteContent: string;
  onChange: (value: string) => void;
  onSave: () => void;
  saving: boolean;
  saved: boolean;
};

export default function NoteDrawer({
  isOpen,
  onClose,
  noteContent,
  onChange,
  onSave,
  saving,
  saved,
}: Props) {
  const textareaRef = useRef<HTMLTextAreaElement>(null);

  // ドロワーが開いたときにテキストエリアにフォーカス
  useEffect(() => {
    if (isOpen) {
      setTimeout(() => textareaRef.current?.focus(), 300);
    }
  }, [isOpen]);

  // ESCキーで閉じる
  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      if (e.key === "Escape" && isOpen) onClose();
    };
    window.addEventListener("keydown", handleKeyDown);
    return () => window.removeEventListener("keydown", handleKeyDown);
  }, [isOpen, onClose]);

  return (
    <>
      {/* オーバーレイ */}
      <div
        className={`fixed inset-0 z-40 bg-black/50 transition-opacity duration-300 ${
          isOpen ? "opacity-100" : "pointer-events-none opacity-0"
        }`}
        onClick={onClose}
      />

      {/* サイドドロワー */}
      <div
        className={`fixed right-0 top-0 z-50 flex h-full w-full max-w-md flex-col bg-white shadow-2xl transition-transform duration-300 ease-in-out ${
          isOpen ? "translate-x-0" : "translate-x-full"
        }`}
      >
        {/* ヘッダー */}
        <div className="flex items-center justify-between border-b border-slate-100 px-6 py-4">
          <h2 className="text-base font-semibold text-slate-800">📝 メモ</h2>
          <button
            onClick={onClose}
            className="flex h-8 w-8 items-center justify-center rounded-full text-slate-400 hover:bg-slate-100 hover:text-slate-600 transition-colors"
            aria-label="閉じる"
          >
            <svg
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 24 24"
              strokeWidth={2}
              stroke="currentColor"
              className="h-5 w-5"
            >
              <path strokeLinecap="round" strokeLinejoin="round" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        {/* コンテンツ */}
        <div className="flex flex-1 flex-col overflow-y-auto px-6 py-5">
          <p className="mb-3 text-xs text-slate-500">
            このレッスンに関するメモを自由に記録できます。
          </p>
          <textarea
            ref={textareaRef}
            value={noteContent}
            onChange={(e) => onChange(e.target.value)}
            placeholder="このレッスンのメモを入力..."
            className="flex-1 min-h-[200px] w-full resize-none rounded-lg border border-slate-200 px-3 py-2 text-sm text-slate-700 focus:border-pink-400 focus:outline-none focus:ring-1 focus:ring-pink-400"
          />
        </div>

        {/* フッター */}
        <div className="border-t border-slate-100 px-6 py-4">
          <div className="flex items-center gap-3">
            <button
              onClick={onSave}
              disabled={saving}
              className="rounded-lg bg-pink-500 px-5 py-2 text-sm font-semibold text-white hover:bg-pink-600 disabled:opacity-50 transition-colors"
            >
              {saving ? "保存中..." : "メモを保存"}
            </button>
            {saved && (
              <span className="flex items-center gap-1 text-sm text-green-600">
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  fill="none"
                  viewBox="0 0 24 24"
                  strokeWidth={2}
                  stroke="currentColor"
                  className="h-4 w-4"
                >
                  <path strokeLinecap="round" strokeLinejoin="round" d="M4.5 12.75l6 6 9-13.5" />
                </svg>
                保存しました
              </span>
            )}
          </div>
        </div>
      </div>
    </>
  );
}
