import Link from "next/link";

const links = {
  Company: ["/about", "/contact", "/portfolio"],
  Services: ["/services", "/packages", "/pricing"],
  Auth: ["/login", "/register", "/forgot-password"],
};

export function SiteFooter() {
  return (
    <footer className="border-t border-slate-200 bg-white">
      <div className="mx-auto grid w-full max-w-7xl gap-10 px-4 py-12 lg:grid-cols-[1.5fr_1fr_1fr_1fr] lg:px-8">
        <div className="space-y-4">
          <div className="flex items-center gap-3">
            <div className="flex h-10 w-10 items-center justify-center rounded-2xl bg-slate-950 text-sm font-bold text-white">NW</div>
            <div>
              <p className="text-sm font-semibold text-slate-900">NexusWeb</p>
              <p className="text-xs text-slate-500">Premium digital services</p>
            </div>
          </div>
          <p className="max-w-md text-sm leading-6 text-slate-600">
            We build websites, applications, and digital products that help businesses move faster and convert better.
          </p>
        </div>

        {Object.entries(links).map(([title, items]) => (
          <div key={title} className="space-y-3">
            <h3 className="text-sm font-semibold text-slate-900">{title}</h3>
            <div className="flex flex-col gap-2 text-sm text-slate-600">
              {items.map((href) => (
                <Link key={href} href={href} className="transition hover:text-slate-950">
                  {href.replace("/", "") || "Home"}
                </Link>
              ))}
            </div>
          </div>
        ))}
      </div>
      <div className="border-t border-slate-200 py-4 text-center text-xs text-slate-500">
        © 2026 NexusWeb Marketplace. All rights reserved.
      </div>
    </footer>
  );
}