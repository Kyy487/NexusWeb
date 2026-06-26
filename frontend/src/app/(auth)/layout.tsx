import { AuthShell } from "@/components/layout/auth-shell";

export default function AuthLayout({ children }: Readonly<{ children: React.ReactNode }>) {
  return <AuthShell title="Welcome back" description="Sign in to manage your projects, invoices, and payments.">{children}</AuthShell>;
}