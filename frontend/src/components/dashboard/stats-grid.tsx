import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";

export function StatsGrid({ stats }: { stats: Array<{ label: string; value: string; note?: string; tone?: "default" | "success" | "warning" | "danger" | "info" }> }) {
  return (
    <div className="grid gap-4 sm:grid-cols-2 xl:grid-cols-4">
      {stats.map((item) => (
        <Card key={item.label}>
          <CardHeader className="pb-2">
            <CardDescription>{item.label}</CardDescription>
            <CardTitle className="text-3xl">{item.value}</CardTitle>
          </CardHeader>
          <CardContent>
            {item.note ? <Badge tone={item.tone ?? "info"}>{item.note}</Badge> : null}
          </CardContent>
        </Card>
      ))}
    </div>
  );
}