import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";

export function AuthShell({
  title,
  description,
  children,
}: {
  title: string;
  description: string;
  children: React.ReactNode;
}) {
  return (
    <div className="min-h-screen bg-[radial-gradient(circle_at_top_left,rgba(37,99,235,0.16),transparent_25%),linear-gradient(135deg,#0f172a_0%,#020617_45%,#f8fafc_45%,#f8fafc_100%)]">
      <div className="mx-auto grid min-h-screen max-w-7xl lg:grid-cols-[1fr_1fr]">
        <div className="flex items-center px-6 py-14 text-white lg:px-10">
          <div className="max-w-xl space-y-6">
            <div className="inline-flex rounded-full border border-white/15 bg-white/10 px-4 py-2 text-xs font-medium uppercase tracking-[0.25em] text-white/80 backdrop-blur">
              NexusWeb Marketplace
            </div>
            <h1 className="text-4xl font-semibold tracking-tight lg:text-6xl">Premium delivery for digital services.</h1>
            <p className="text-base leading-7 text-white/70">
              A modern client portal designed to feel trustworthy, fast, and built for conversion.
            </p>
            <div className="grid gap-3 sm:grid-cols-2">
              {[
                "JWT Authentication",
                "Protected dashboards",
                "Order to payment flow",
                "Responsive premium UI",
              ].map((item) => (
                <Card key={item} className="border-white/10 bg-white/5 text-white shadow-none backdrop-blur">
                  <CardContent className="p-4 text-sm text-white/80">{item}</CardContent>
                </Card>
              ))}
            </div>
          </div>
        </div>

        <div className="flex items-center justify-center px-6 py-14 lg:px-10">
          <Card className="w-full max-w-md border-slate-200/80 bg-white/95 shadow-[0_24px_80px_-35px_rgba(15,23,42,0.45)]">
            <CardHeader>
              <CardTitle className="text-2xl">{title}</CardTitle>
              <CardDescription>{description}</CardDescription>
            </CardHeader>
            <CardContent>{children}</CardContent>
          </Card>
        </div>
      </div>
    </div>
  );
}