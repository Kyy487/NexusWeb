import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { cn } from "@/lib/utils";

export function Timeline({
  steps,
  currentIndex = 0,
}: Readonly<{
  steps: Array<{ title: string; description: string }>;
  currentIndex?: number;
}>) {
  return (
    <Card>
      <CardHeader>
        <CardTitle className="text-lg">Project Timeline</CardTitle>
        <CardDescription>Requirement → Invoice → Payment → In Progress → Revision → Completed</CardDescription>
      </CardHeader>
      <CardContent className="space-y-5">
        {steps.map((step, index) => {
          const active = index <= currentIndex;
          return (
            <div key={step.title} className="flex gap-4">
              <div className="flex flex-col items-center">
                <div className={cn("h-4 w-4 rounded-full border-4", active ? "border-blue-600 bg-blue-600" : "border-slate-300 bg-white")} />
                {index < steps.length - 1 ? <div className={cn("mt-1 h-full w-px flex-1", active ? "bg-blue-200" : "bg-slate-200")} /> : null}
              </div>
              <div className="pb-5">
                <h4 className={cn("font-semibold", active ? "text-slate-950" : "text-slate-500")}>{step.title}</h4>
                <p className="text-sm leading-6 text-slate-600">{step.description}</p>
              </div>
            </div>
          );
        })}
      </CardContent>
    </Card>
  );
}