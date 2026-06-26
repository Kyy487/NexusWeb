import { DashboardShell } from "@/components/layout/dashboard-shell";

const items = [
  { href: "/customer/dashboard", label: "Dashboard" },
  { href: "/customer/orders", label: "Orders" },
  { href: "/customer/invoices", label: "Invoices" },
  { href: "/customer/payments", label: "Payments" },
  { href: "/customer/projects", label: "Projects" },
  { href: "/customer/files", label: "Files" },
  { href: "/customer/messages", label: "Messages" },
  { href: "/customer/profile", label: "Profile" },
  { href: "/customer/settings", label: "Settings" },
];

export default function CustomerLayout({ children }: Readonly<{ children: React.ReactNode }>) {
  return <DashboardShell title="Customer Portal" items={items}>{children}</DashboardShell>;
}