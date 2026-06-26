"use client";

import Link from "next/link";
import { usePathname, useRouter } from "next/navigation";

import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { useAuthStore } from "@/store/auth-store";
import { cn } from "@/lib/utils";

type NavItem = { href: string; label: string };

function DashboardSidebar({ title, items }: { title: string; items: NavItem[] }) {
  const pathname = usePathname();
  return (
    <aside className="hidden min-h-screen border-r border-slate-200 bg-white/90 px-4 py-6 lg:block lg:w-72">
      <div className="flex items-center gap-3 px-2 pb-8">
        <div className="flex h-10 w-10 items-center justify-center rounded-2xl bg-slate-950 text-sm font-bold text-white">NW</div>
        <div>
          <p className="text-sm font-semibold text-slate-900">NexusWeb</p>
          <p className="text-xs text-slate-500">{title}</p>
        </div>
      </div>
      <nav className="space-y-1">
        {items.map((item) => {
          const active = pathname === item.href || pathname.startsWith(`${item.href}/`);
          return (
            <Link
              key={item.href}
              href={item.href}
              className={cn(
                "flex items-center justify-between rounded-2xl px-4 py-3 text-sm transition",
                active ? "bg-slate-950 text-white" : "text-slate-600 hover:bg-slate-100 hover:text-slate-950"
              )}
            >
              {item.label}
              {active ? <Badge tone="info">Active</Badge> : null}
            </Link>
          );
        })}
      </nav>
    </aside>
  );
}

export function DashboardShell({
  title,
  items,
  children,
}: {
  title: string;
  items: NavItem[];
  children: React.ReactNode;
}) {
  const router = useRouter();
  const clearAuth = useAuthStore((state) => state.clearAuth);
  const user = useAuthStore((state) => state.user);

  const handleLogout = () => {
    clearAuth();
    router.push("/login");
  };

  return (
    <div className="min-h-screen bg-slate-50 lg:flex">
      <DashboardSidebar title={title} items={items} />
      <div className="flex-1">
        <header className="sticky top-0 z-30 border-b border-slate-200/80 bg-white/85 backdrop-blur-xl">
          <div className="flex items-center justify-between gap-4 px-4 py-4 lg:px-8">
            <div>
              <p className="text-sm text-slate-500">Welcome back</p>
              <h1 className="text-xl font-semibold text-slate-900">{user?.name ?? title}</h1>
            </div>
            <div className="flex items-center gap-3">
              {user ? <Badge tone="info">{user.role}</Badge> : null}
              <Button variant="outline" onClick={handleLogout}>
                Logout
              </Button>
            </div>
          </div>
        </header>
        <main className="px-4 py-6 lg:px-8 lg:py-8">{children}</main>
      </div>
    </div>
  );
}