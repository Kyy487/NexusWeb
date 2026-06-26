import Link from "next/link";

import { Button } from "@/components/ui/button";

const navItems = [
  { href: "/services", label: "Services" },
  { href: "/packages", label: "Packages" },
  { href: "/pricing", label: "Pricing" },
  { href: "/portfolio", label: "Portfolio" },
  { href: "/about", label: "About" },
  { href: "/contact", label: "Contact" },
];

export function SiteHeader() {
  return (
    <header className="sticky top-0 z-40 border-b border-slate-200/70 bg-white/80 backdrop-blur-xl">
      <div className="mx-auto flex w-full max-w-7xl items-center justify-between px-4 py-4 lg:px-8">
        <Link href="/" className="flex items-center gap-3">
          <div className="flex h-10 w-10 items-center justify-center rounded-2xl bg-slate-950 text-sm font-bold text-white">NW</div>
          <div>
            <p className="text-sm font-semibold text-slate-900">NexusWeb</p>
            <p className="text-xs text-slate-500">Digital Services Marketplace</p>
          </div>
        </Link>

        <nav className="hidden items-center gap-6 text-sm text-slate-600 lg:flex">
          {navItems.map((item) => (
            <Link key={item.href} href={item.href} className="transition hover:text-slate-950">
              {item.label}
            </Link>
          ))}
        </nav>

        <div className="flex items-center gap-3">
          <Button variant="ghost" asChild>
            <Link href="/login">Login</Link>
          </Button>
          <Button asChild>
            <Link href="/register">Start Project</Link>
          </Button>
        </div>
      </div>
    </header>
  );
}