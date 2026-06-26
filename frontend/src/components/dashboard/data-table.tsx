import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";

export function DataTable<T extends { id: string }>({
  title,
  description,
  columns,
  rows,
}: Readonly<{
  title: string;
  description: string;
  columns: string[];
  rows: T[];
}>) {
  return (
    <Card>
      <CardHeader>
        <CardTitle className="text-lg">{title}</CardTitle>
        <CardDescription>{description}</CardDescription>
      </CardHeader>
      <CardContent>
        <div className="overflow-x-auto">
          <table className="min-w-full text-left text-sm">
            <thead>
              <tr className="border-b border-slate-200 text-slate-500">
                {columns.map((column) => (
                  <th key={column} className="py-3 pr-4 font-medium">
                    {column}
                  </th>
                ))}
              </tr>
            </thead>
            <tbody>
              {rows.map((row) => (
                <tr key={row.id} className="border-b border-slate-100 last:border-0">
                  {columns.map((column) => {
                    const normalizedKey = column.toLowerCase().replace(/\s+/g, "_");
                    const candidateKeys = [normalizedKey, column, column.toLowerCase()];
                    const rowRecord = row as Record<string, unknown>;
                    const value = candidateKeys.map((key) => rowRecord[key]).find((candidate) => candidate !== undefined);

                    return (
                      <td key={column} className="py-3 pr-4 text-slate-700">
                        {value == null ? "-" : typeof value === "object" ? JSON.stringify(value) : String(value)}
                      </td>
                    );
                  })}
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </CardContent>
    </Card>
  );
}