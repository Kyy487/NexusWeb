import Link from "next/link";
import { ArrowRight, BadgeCheck, Globe, Layers3, PenTool, Rocket, ShieldCheck, Star } from "lucide-react";

import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Separator } from "@/components/ui/separator";
import { cn } from "@/lib/utils";

const services = [
  { icon: Globe, title: "Website Development", description: "Conversion-focused websites built for speed, clarity, and trust." },
  { icon: Layers3, title: "Web Application", description: "Modern product interfaces for internal tools and customer platforms." },
  { icon: PenTool, title: "UI/UX Design", description: "Premium interface systems that make digital products feel expensive." },
  { icon: Rocket, title: "Launch Support", description: "Delivery, testing, and launch readiness for production shipping." },
];

const packages = [
  { name: "Starter", price: "From Rp 3.5M", features: ["Landing page", "Responsive UI", "Basic SEO"] },
  { name: "Growth", price: "From Rp 9.5M", features: ["Multi-page website", "CMS-ready", "Premium UI system"] },
  { name: "Pro", price: "Custom", features: ["Web app", "Dashboard", "Priority support"] },
];

const testimonials = [
  { name: "Ari Wibowo", role: "Founder, SaaS Studio", text: "The team delivered a sharp, premium portal that improved lead conversion immediately." },
  { name: "Nadia Putri", role: "Operations Lead", text: "Fast communication, clean execution, and a workflow that felt production-ready from day one." },
  { name: "Kevin Santoso", role: "Marketing Manager", text: "The final experience looked like a modern SaaS product, not a template site." },
];

const faqs = [
  ["How do I start?", "Create an account, submit a project order, and we will guide the next steps through invoice and payment."],
  ["Can I track my project?", "Yes. Customer project pages expose order status, invoice, payment, files, and progress timeline."],
  ["Is payment secure?", "Payment is handled through the existing invoice and payment flow with verification steps."],
];

function SectionTitle({ eyebrow, title, description }: Readonly<{ eyebrow: string; title: string; description: string }>) {
  return (
    <div className="mx-auto max-w-2xl text-center">
      <Badge tone="info" className="mb-4">
        {eyebrow}
      </Badge>
      <h2 className="text-3xl font-semibold tracking-tight text-slate-950 lg:text-4xl">{title}</h2>
      <p className="mt-4 text-base leading-7 text-slate-600">{description}</p>
    </div>
  );
}

function StatPill({ label }: Readonly<{ label: string }>) {
  return (
    <div className="flex items-center gap-2 rounded-full border border-white/60 bg-white/85 px-4 py-2 text-sm font-medium text-slate-700 shadow-lg shadow-slate-900/5 backdrop-blur">
      <BadgeCheck className="h-4 w-4 text-emerald-500" />
      {label}
    </div>
  );
}

export function LandingPage() {
  return (
    <div>
      <section className="mx-auto grid max-w-7xl gap-12 px-4 pb-20 pt-16 lg:grid-cols-[1.1fr_0.9fr] lg:px-8 lg:py-24">
        <div className="space-y-8">
          <Badge tone="info">Premium Digital Services Marketplace</Badge>
          <div className="space-y-5">
            <h1 className="max-w-3xl text-5xl font-semibold tracking-tight text-slate-950 lg:text-7xl">
              Build Professional Digital Solutions for Your Business.
            </h1>
            <p className="max-w-2xl text-lg leading-8 text-slate-600">
              We help businesses grow through websites, applications, and modern digital services with a premium client experience.
            </p>
          </div>
          <div className="flex flex-wrap gap-3">
            <Button asChild size="lg">
              <Link href="/register">Start Project <ArrowRight className="ml-2 h-4 w-4" /></Link>
            </Button>
            <Button asChild variant="outline" size="lg">
              <Link href="/services">Explore Services</Link>
            </Button>
          </div>
          <div className="grid gap-3 sm:grid-cols-2 xl:grid-cols-4">
            {["100+ Projects Delivered", "Professional Team", "Transparent Process", "Secure Payment"].map((item) => (
              <StatPill key={item} label={item} />
            ))}
          </div>
        </div>

        <Card className="relative overflow-hidden border-slate-200/80 bg-slate-950 text-white">
          <div className="absolute inset-0 bg-[radial-gradient(circle_at_top_right,rgba(37,99,235,0.3),transparent_35%),radial-gradient(circle_at_bottom_left,rgba(6,182,212,0.18),transparent_28%)]" />
          <CardHeader className="relative">
            <CardTitle className="text-white">Client Project Snapshot</CardTitle>
            <CardDescription className="text-slate-300">A polished journey from order to completion.</CardDescription>
          </CardHeader>
          <CardContent className="relative space-y-4">
            {[
              ["Order Created", "Web redesign for a growing startup"],
              ["Invoice Issued", "Awaiting secure payment verification"],
              ["Project In Progress", "Design, build, and review cycles"],
            ].map(([title, desc], index) => (
              <div key={title} className="flex items-start gap-4 rounded-3xl border border-white/10 bg-white/5 p-4">
                <div className="mt-1 flex h-8 w-8 items-center justify-center rounded-full bg-blue-500 text-sm font-semibold text-white">{index + 1}</div>
                <div>
                  <p className="font-semibold text-white">{title}</p>
                  <p className="text-sm text-slate-300">{desc}</p>
                </div>
              </div>
            ))}
            <Separator className="bg-white/10" />
            <div className="flex flex-wrap gap-3 text-xs text-slate-300">
              <StatPill label="Project Timeline" />
              <StatPill label="Invoice Tracking" />
              <StatPill label="File Uploads" />
            </div>
          </CardContent>
        </Card>
      </section>

      <section className="mx-auto max-w-7xl px-4 py-16 lg:px-8">
        <SectionTitle eyebrow="Trusted By" title="Built for premium client experiences" description="NexusWeb is designed to feel like a modern SaaS storefront while staying focused on delivery, trust, and conversion." />
      </section>

      <section className="mx-auto max-w-7xl px-4 py-16 lg:px-8">
        <SectionTitle eyebrow="Services" title="Everything you need to launch and grow" description="High-conversion design, development, and delivery support across the project lifecycle." />
        <div className="mt-10 grid gap-6 md:grid-cols-2 xl:grid-cols-4">
          {services.map((service) => {
            const Icon = service.icon;
            return (
              <Card key={service.title} className="group transition hover:-translate-y-1 hover:shadow-[0_30px_80px_-35px_rgba(15,23,42,0.35)]">
                <CardHeader>
                  <div className="mb-5 flex h-12 w-12 items-center justify-center rounded-2xl bg-blue-50 text-blue-700 transition group-hover:bg-blue-600 group-hover:text-white">
                    <Icon className="h-5 w-5" />
                  </div>
                  <CardTitle>{service.title}</CardTitle>
                  <CardDescription>{service.description}</CardDescription>
                </CardHeader>
              </Card>
            );
          })}
        </div>
      </section>

      <section className="mx-auto max-w-7xl px-4 py-16 lg:px-8">
        <SectionTitle eyebrow="How It Works" title="A simple workflow that keeps the project moving" description="Register, place the order, submit requirements, receive invoice, pay, and track delivery in one place." />
        <div className="mt-10 grid gap-6 lg:grid-cols-5">
          {[
            ["Register", "Create your account and enter the platform."],
            ["Create Order", "Choose service and package, then submit project details."],
            ["Requirements", "Add the project brief and supporting files."],
            ["Invoice & Payment", "Review invoice and submit payment proof."],
            ["Track Project", "Follow progress until completion."],
          ].map(([title, description], index) => (
            <Card key={title} className="relative overflow-hidden">
              <CardContent className="space-y-4">
                <div className="flex h-10 w-10 items-center justify-center rounded-2xl bg-slate-950 text-sm font-semibold text-white">0{index + 1}</div>
                <h3 className="text-lg font-semibold text-slate-950">{title}</h3>
                <p className="text-sm leading-6 text-slate-600">{description}</p>
              </CardContent>
            </Card>
          ))}
        </div>
      </section>

      <section className="mx-auto max-w-7xl px-4 py-16 lg:px-8">
        <SectionTitle eyebrow="Featured Packages" title="High-value offers tailored to different growth stages" description="Packages are presented to help buyers quickly understand scope, value, and next steps." />
        <div className="mt-10 grid gap-6 lg:grid-cols-3">
          {packages.map((item, index) => (
            <Card key={item.name} className={cn(index === 1 ? "border-blue-200 shadow-[0_30px_80px_-35px_rgba(37,99,235,0.28)]" : "") }>
              <CardHeader>
                <CardTitle>{item.name}</CardTitle>
                <CardDescription>{item.price}</CardDescription>
              </CardHeader>
              <CardContent>
                <ul className="space-y-3 text-sm text-slate-600">
                  {item.features.map((feature) => (
                    <li key={feature} className="flex items-center gap-2">
                      <ShieldCheck className="h-4 w-4 text-emerald-500" /> {feature}
                    </li>
                  ))}
                </ul>
              </CardContent>
            </Card>
          ))}
        </div>
      </section>

      <section className="mx-auto max-w-7xl px-4 py-16 lg:px-8">
        <SectionTitle eyebrow="Portfolio" title="Designed to feel premium and trustworthy" description="A curated showcase to communicate quality before the first message is sent." />
        <div className="mt-10 grid gap-6 md:grid-cols-2 xl:grid-cols-3">
          {[
            "Corporate Website Redesign",
            "Digital Agency Portal",
            "Customer Operations Dashboard",
          ].map((item) => (
            <Card key={item} className="overflow-hidden">
              <div className="h-56 bg-[linear-gradient(135deg,#0f172a,#1e3a8a_45%,#06b6d4)]" />
              <CardContent>
                <h3 className="text-lg font-semibold text-slate-950">{item}</h3>
                <p className="mt-2 text-sm text-slate-600">Premium visual system, responsive layout, and a conversion-oriented journey.</p>
              </CardContent>
            </Card>
          ))}
        </div>
      </section>

      <section className="mx-auto max-w-7xl px-4 py-16 lg:px-8">
        <SectionTitle eyebrow="Testimonials" title="Trusted by teams that care about quality" description="The product should feel serious, modern, and easy to trust." />
        <div className="mt-10 grid gap-6 lg:grid-cols-3">
          {testimonials.map((item) => (
            <Card key={item.name}>
              <CardHeader>
                <div className="mb-4 flex items-center gap-1 text-amber-400">
                  {Array.from({ length: 5 }).map((_, index) => (
                    <Star key={`star-${index}`} className="h-4 w-4 fill-current" />
                  ))}
                </div>
                <CardDescription className="text-base text-slate-700">{item.text}</CardDescription>
              </CardHeader>
              <CardContent>
                <p className="font-semibold text-slate-950">{item.name}</p>
                <p className="text-sm text-slate-500">{item.role}</p>
              </CardContent>
            </Card>
          ))}
        </div>
      </section>

      <section className="mx-auto max-w-4xl px-4 py-16 lg:px-8">
        <SectionTitle eyebrow="FAQ" title="Common questions before starting a project" description="A clear project flow reduces friction and keeps conversion high." />
        <div className="mt-10 space-y-4">
          {faqs.map(([question, answer]) => (
            <Card key={question}>
              <CardHeader>
                <CardTitle className="text-base">{question}</CardTitle>
                <CardDescription>{answer}</CardDescription>
              </CardHeader>
            </Card>
          ))}
        </div>
      </section>

      <section className="mx-auto max-w-7xl px-4 py-20 lg:px-8">
        <Card className="overflow-hidden border-slate-200 bg-slate-950 text-white">
          <div className="grid gap-10 p-8 lg:grid-cols-[1.2fr_0.8fr] lg:p-12">
            <div className="space-y-4">
              <Badge tone="info">Final CTA</Badge>
              <h2 className="text-3xl font-semibold tracking-tight lg:text-5xl">Ready to build something that feels premium?</h2>
              <p className="max-w-2xl text-slate-300">Start a project, move through the invoice and payment flow, and track delivery in a polished customer experience.</p>
            </div>
            <div className="flex flex-col justify-center gap-3 lg:items-end">
              <Button asChild size="lg">
                <Link href="/register">Start Project</Link>
              </Button>
              <Button asChild variant="outline" size="lg" className="border-white/15 bg-white/5 text-white hover:bg-white/10">
                <Link href="/login">Login to Dashboard</Link>
              </Button>
            </div>
          </div>
        </Card>
      </section>
    </div>
  );
}