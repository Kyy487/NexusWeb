"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { useMutation } from "@tanstack/react-query";
import { Loader2 } from "lucide-react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { useForm } from "react-hook-form";
import { z } from "zod";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { authApi } from "@/lib/api";
import { useAuthStore } from "@/store/auth-store";

const loginSchema = z.object({
  email: z.string().email("Enter a valid email"),
  password: z.string().min(6, "Password must be at least 6 characters"),
});

const registerSchema = loginSchema.extend({
  name: z.string().min(2, "Name is required"),
  phone: z.string().optional(),
});

const forgotSchema = z.object({ email: z.string().email("Enter a valid email") });

type LoginValues = z.infer<typeof loginSchema>;
type RegisterValues = z.infer<typeof registerSchema>;

export function LoginForm() {
  const router = useRouter();
  const setAuth = useAuthStore((state) => state.setAuth);
  const { register, handleSubmit, formState: { errors } } = useForm<LoginValues>({ resolver: zodResolver(loginSchema) });

  const mutation = useMutation({
    mutationFn: authApi.login,
    onSuccess: (data) => {
      setAuth(data.token, data.user);
      router.push(data.user.role === "CUSTOMER" ? "/customer/dashboard" : "/admin/dashboard");
    },
  });

  return (
    <form className="space-y-4" onSubmit={handleSubmit((values) => mutation.mutate(values))}>
      <div>
        <label htmlFor="login-email" className="mb-2 block text-sm font-medium text-slate-700">Email</label>
        <Input id="login-email" type="email" placeholder="name@company.com" {...register("email")} />
        {errors.email ? <p className="mt-2 text-xs text-red-500">{errors.email.message}</p> : null}
      </div>
      <div>
        <label htmlFor="login-password" className="mb-2 block text-sm font-medium text-slate-700">Password</label>
        <Input id="login-password" type="password" placeholder="••••••••" {...register("password")} />
        {errors.password ? <p className="mt-2 text-xs text-red-500">{errors.password.message}</p> : null}
      </div>
      <Button className="w-full" type="submit" disabled={mutation.isPending}>
        {mutation.isPending ? <Loader2 className="mr-2 h-4 w-4 animate-spin" /> : null} Login
      </Button>
      <div className="flex items-center justify-between text-sm text-slate-600">
        <Link href="/forgot-password" className="transition hover:text-slate-950">Forgot password?</Link>
        <Link href="/register" className="transition hover:text-slate-950">Create account</Link>
      </div>
      {mutation.error ? <p className="rounded-2xl bg-red-50 px-4 py-3 text-sm text-red-700">Unable to login. Check credentials or backend connection.</p> : null}
    </form>
  );
}

export function RegisterForm() {
  const router = useRouter();
  const setAuth = useAuthStore((state) => state.setAuth);
  const { register, handleSubmit, formState: { errors } } = useForm<RegisterValues>({ resolver: zodResolver(registerSchema) });

  const mutation = useMutation({
    mutationFn: authApi.register,
    onSuccess: (data) => {
      setAuth(data.token, data.user);
      router.push("/customer/dashboard");
    },
  });

  return (
    <form className="space-y-4" onSubmit={handleSubmit((values) => mutation.mutate(values))}>
      <div>
        <label htmlFor="register-name" className="mb-2 block text-sm font-medium text-slate-700">Full name</label>
        <Input id="register-name" placeholder="Your name" {...register("name")} />
        {errors.name ? <p className="mt-2 text-xs text-red-500">{errors.name.message}</p> : null}
      </div>
      <div className="grid gap-4 sm:grid-cols-2">
        <div>
          <label htmlFor="register-email" className="mb-2 block text-sm font-medium text-slate-700">Email</label>
          <Input id="register-email" type="email" placeholder="name@company.com" {...register("email")} />
          {errors.email ? <p className="mt-2 text-xs text-red-500">{errors.email.message}</p> : null}
        </div>
        <div>
          <label htmlFor="register-phone" className="mb-2 block text-sm font-medium text-slate-700">Phone</label>
          <Input id="register-phone" placeholder="08xxxx" {...register("phone")} />
        </div>
      </div>
      <div>
        <label htmlFor="register-password" className="mb-2 block text-sm font-medium text-slate-700">Password</label>
        <Input id="register-password" type="password" placeholder="••••••••" {...register("password")} />
        {errors.password ? <p className="mt-2 text-xs text-red-500">{errors.password.message}</p> : null}
      </div>
      <Button className="w-full" type="submit" disabled={mutation.isPending}>
        {mutation.isPending ? <Loader2 className="mr-2 h-4 w-4 animate-spin" /> : null} Create account
      </Button>
      {mutation.error ? <p className="rounded-2xl bg-red-50 px-4 py-3 text-sm text-red-700">Unable to register. Check backend availability.</p> : null}
    </form>
  );
}

export function ForgotPasswordForm() {
  const { register, formState: { errors } } = useForm<{ email: string }>({ resolver: zodResolver(forgotSchema) });

  return (
    <form className="space-y-4">
      <div>
        <label htmlFor="forgot-email" className="mb-2 block text-sm font-medium text-slate-700">Email</label>
        <Input id="forgot-email" type="email" placeholder="name@company.com" {...register("email")} />
        {errors.email ? <p className="mt-2 text-xs text-red-500">{errors.email.message}</p> : null}
      </div>
      <Button className="w-full" type="submit">Send reset instructions</Button>
      <p className="text-sm text-slate-500">Password reset UI is ready for backend integration.</p>
    </form>
  );
}

const orderSchema = z.object({
  customer_id: z.string().min(1),
  service_id: z.string().min(1),
  package_id: z.string().min(1),
  title: z.string().min(2),
  description: z.string().optional(),
  deadline: z.string().optional(),
});

export function OrderForm({
  services,
  packages,
  onSubmit,
}: Readonly<{
  services: Array<{ id: string; name: string; category_name?: string }>;
  packages: Array<{ id: string; name: string; service_id: string; price: number; revision_count: number; delivery_days: number }>;
  onSubmit: (values: z.infer<typeof orderSchema>) => void;
}>) {
  const user = useAuthStore((state) => state.user);
  const { register, handleSubmit, watch } = useForm<z.infer<typeof orderSchema>>({
    resolver: zodResolver(orderSchema),
    defaultValues: { customer_id: user?.id ?? "" },
  });

  const serviceId = watch("service_id");
  const filteredPackages = packages.filter((item) => !serviceId || item.service_id === serviceId);

  return (
    <form className="space-y-6" onSubmit={handleSubmit(onSubmit)}>
      <input type="hidden" {...register("customer_id")} />
      <div className="grid gap-4 lg:grid-cols-2">
        <div>
          <label htmlFor="service_id" className="mb-2 block text-sm font-medium text-slate-700">Choose Service</label>
          <select id="service_id" {...register("service_id")} className="h-11 w-full rounded-2xl border border-slate-200 bg-white px-4 text-sm outline-none focus:border-blue-500">
            <option value="">Select service</option>
            {services.map((service) => (
              <option key={service.id} value={service.id}>{service.name}</option>
            ))}
          </select>
        </div>
        <div>
          <label htmlFor="package_id" className="mb-2 block text-sm font-medium text-slate-700">Choose Package</label>
          <select id="package_id" {...register("package_id")} className="h-11 w-full rounded-2xl border border-slate-200 bg-white px-4 text-sm outline-none focus:border-blue-500">
            <option value="">Select package</option>
            {filteredPackages.map((item) => (
              <option key={item.id} value={item.id}>{item.name} - Rp {item.price.toLocaleString("id-ID")}</option>
            ))}
          </select>
        </div>
      </div>
      <div className="grid gap-4 lg:grid-cols-2">
        <div>
          <label htmlFor="project-title" className="mb-2 block text-sm font-medium text-slate-700">Project Title</label>
          <Input id="project-title" placeholder="Company profile website" {...register("title")} />
        </div>
        <div>
          <label htmlFor="deadline" className="mb-2 block text-sm font-medium text-slate-700">Deadline</label>
          <Input id="deadline" type="date" {...register("deadline")} />
        </div>
      </div>
      <div>
        <label htmlFor="description" className="mb-2 block text-sm font-medium text-slate-700">Description</label>
        <Textarea id="description" placeholder="Explain your project scope and goals" {...register("description")} />
      </div>
      <Button type="submit">Review and Submit</Button>
    </form>
  );
}

export function RequirementForm({ onSubmit }: Readonly<{ onSubmit: (values: { question: string; answer: string }) => void }>) {
  const { register, handleSubmit } = useForm<{ question: string; answer: string }>();
  return (
    <form className="space-y-4" onSubmit={handleSubmit(onSubmit)}>
      <div>
        <label htmlFor="question" className="mb-2 block text-sm font-medium text-slate-700">Question</label>
        <Textarea id="question" placeholder='"What content and sections are required?"' {...register("question")} />
      </div>
      <div>
        <label htmlFor="answer" className="mb-2 block text-sm font-medium text-slate-700">Answer</label>
        <Textarea id="answer" placeholder="Provide the actual requirement detail" {...register("answer")} />
      </div>
      <Button type="submit">Save Requirement</Button>
    </form>
  );
}

export function PaymentForm({ onSubmit }: Readonly<{ onSubmit: (values: { amount: number; payment_method?: string; payment_proof_url?: string }) => void }>) {
  const { register, handleSubmit } = useForm<{ amount: number; payment_method?: string; payment_proof_url?: string }>();
  return (
    <form className="space-y-4" onSubmit={handleSubmit(onSubmit)}>
      <div className="grid gap-4 lg:grid-cols-2">
        <div>
          <label htmlFor="amount" className="mb-2 block text-sm font-medium text-slate-700">Amount</label>
          <Input id="amount" type="number" step="0.01" {...register("amount", { valueAsNumber: true })} />
        </div>
        <div>
          <label htmlFor="payment-method" className="mb-2 block text-sm font-medium text-slate-700">Payment Method</label>
          <Input id="payment-method" placeholder="Bank Transfer / QRIS" {...register("payment_method")} />
        </div>
      </div>
      <div>
        <label htmlFor="payment-proof-url" className="mb-2 block text-sm font-medium text-slate-700">Payment Proof URL</label>
        <Input id="payment-proof-url" placeholder="https://..." {...register("payment_proof_url")} />
      </div>
      <Button type="submit">Submit Payment</Button>
    </form>
  );
}