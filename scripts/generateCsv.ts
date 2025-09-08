import PocketBase from 'pocketbase';
import fs from 'fs';
import path from 'path';
import { fileURLToPath } from 'url';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

const OUTPUT_PATH = path.join(__dirname, '..', 'public', 'suggestions.csv');

(async () => {
  const pb = new PocketBase('http://127.0.0.1:8090');

  const suggestionTypes = await pb.collection("SuggestionType").getFullList({
    filter: "isActive=true",
    sort: "columnOrder",
  });

  const suggestions = await pb.collection("Suggestions").getFullList({
    filter: "active=true",
  });

  const columnsMeta: [string, string][] = suggestionTypes.map((t: any) => [t.id, t.name]);
  const columnNames = columnsMeta.map(([, name]) => name);

  const grouped = new Map<string, any[]>();
  for (const s of suggestions) {
    if (!grouped.has(s.idUniqueID)) grouped.set(s.idUniqueID, []);
    grouped.get(s.idUniqueID)!.push(s);
  }

  const rows: Record<string, string>[] = Array.from(grouped.entries()).map(([uniqueId, records]) => {
    const row: Record<string, string> = { idUniqueID: uniqueId };

    for (const [typeId, name] of columnsMeta) {
      const candidates = records.filter((r) => r.idSuggestionType === typeId);
      let val = "";

      if (candidates.length) {
        candidates.sort((a, b) => b.confidence - a.confidence);
        val = candidates[0].suggestion;
      }

      row[name] = val?.toString().trim() || "";
    }

    return row;
  });

  const header = ['idUniqueID', ...columnNames];
  const csv = [
    header.join(','),
    ...rows.map(row =>
      header.map(col => `"${(row[col] || '').replace(/"/g, '""')}"`).join(',')
    ),
  ].join('\n');

  fs.writeFileSync(OUTPUT_PATH, csv, 'utf8');
  console.log(`Wrote ${rows.length} rows to ${OUTPUT_PATH}`);
})();
