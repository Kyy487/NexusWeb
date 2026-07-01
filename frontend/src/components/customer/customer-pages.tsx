"use client";

import { useMutation, useQuery } from "@tanstack/react-query";
import { useRouter } from "next/navigation";
import { useMemo, useState } from "react";
import { Loader2, Upload } from "lucide-react";

import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { DataTable } from "@/components/dashboard/data-table";
import { EmptyState } from "@/components/dashboard/empty-state";
import { PageHeader } from "@/components/dashboard/page-header";
import { StatsGrid } from "@/components/dashboard/stats-grid";
import { Timeline } from "@/components/dashboard/timeline";
import { OrderForm, RequirementForm } from "@/components/forms/auth-forms";
import {
  catalogApi,
  dashboardApi,
  fileApi,
  invoiceApi,
  orderApi,
  paymentApi,
  progressApi,
  requirementApi,
} from "@/lib/api";
import { useAuthStore } from "@/store/auth-store";
import { useProjectStore } from "@/store/project-store";
import type { ActivityLogItem, InvoiceItem, OrderItem, PaymentItem, ProgressItem } from "@/types/api";

function currency(value?: number | null) {
  return new Intl.NumberFormat("id-ID", {
    style: "currency",
    currency: "IDR",
    maximumFractionDigits: 0,
  }).format(value ?? 0);
}

function statusTone(status?: string) {
  if (!status) return "default";
  if (["PAID", "COMPLETED", "ACTIVE", "IN_PROGRESS"].includes(status))
    return "success";
  if (["PENDING", "UNPAID", "OVERDUE", "REVISION"].includes(status))
    return "warning";
  if (["CANCELLED", "FAILED", "EXPIRED"].includes(status)) return "danger";
  return "info";
}

function statusIndex(status?: string) {
  switch (status) {
    case "PENDING":
    case "UNPAID":
      return 0;
    case "PAID":
      return 2;
    case "IN_PROGRESS":
      return 3;
    case "REVISION":
      return 4;
    case "COMPLETED":
      return 5;
    default:
      return 1;
  }
}

export function CustomerDashboardView() {
  const user = useAuthStore((state) => state.user);
  const setCurrentOrder = useProjectStore((state) => state.setCurrentOrder);

  // Fetch real data from APIs
  const {
    data: stats,
    isLoading: statsLoading,
    isError: statsError,
    refetch: refetchStats,
  } = useQuery({
    queryKey: ["customer-stats"],
    queryFn: dashboardApi.customerStats,
  });

  const {
    data: orders = [],
    isLoading: ordersLoading,
    isError: ordersError,
    refetch: refetchOrders,
  } = useQuery({
    queryKey: ["my-orders"],
    queryFn: orderApi.myOrders,
  });

  const {
    data: invoices = [],
    isLoading: invoicesLoading,
  } = useQuery({
    queryKey: ["my-invoices"],
    queryFn: invoiceApi.myInvoices,
  });

  const {
    data: activities = [],
    isLoading: activitiesLoading,
  } = useQuery({
    queryKey: ["my-activities"],
    queryFn: dashboardApi.activities,
  });

  // Get the most recent order
  const latestOrder = orders[0] ?? null;

  // Get progress for latest order if any
  const {
    data: progress = [],
    isLoading: progressLoading,
  } = useQuery({
    queryKey: ["order-progress", latestOrder?.id],
    queryFn: () => progressApi.byOrderId(latestOrder!.id),
    enabled: !!latestOrder?.id,
  });

  const statsCards = stats
    ? [
        {
          label: "Total Orders",
          value: stats.total_orders.toString(),
          note: "All time",
          tone: "info" as const,
        },
        {
          label: "Active Projects",
          value: stats.active_projects.toString(),
          note: "In progress",
          tone: "warning" as const,
        },
        {
          label: "Pending Invoices",
          value: stats.pending_invoices.toString(),
          note: "Awaiting payment",
          tone: "danger" as const,
        },
        {
          label: "Completed",
          value: stats.completed_projects.toString(),
          note: "Delivered",
          tone: "success" as const,
        },
      ]
    : [];

  // Latest pending invoice
  const pendingInvoice = invoices.find(
    (i) => i.status === "UNPAID" || i.status === "OVERDUE"
  );

  return (
    <div className="space-y-8">
      <PageHeader
        eyebrow="Customer Dashboard"
        title={`Welcome back${user?.name ? `, ${user.name}` : ""}!`}
        description="Track your project lifecycle from order to completion in one unified workspace."
      />

      {/* Stats Grid */}
      {statsLoading ? (
        <div className="grid gap-4 sm:grid-cols-2 xl:grid-cols-4">
          {[1, 2, 3, 4].map((n) => (
            <Card key={n}>
              <CardHeader className="pb-2">
                <div className="h-4 w-24 animate-pulse rounded bg-slate-200" />
                <div className="mt-2 h-8 w-16 animate-pulse rounded bg-slate-200" />
              </CardHeader>
              <CardContent>
                <div className="h-5 w-20 animate-pulse rounded bg-slate-100" />
              </CardContent>
            </Card>
          ))}
        </div>
      ) : statsError ? (
        <Card className="border-red-100 bg-red-50">
          <CardContent className="flex items-center justify-between py-4">
            <p className="text-sm text-red-700">Failed to load dashboard stats.</p>
            <Button size="sm" variant="outline" onClick={() => refetchStats()}>
              Retry
            </Button>
          </CardContent>
        </Card>
      ) : (
        <StatsGrid stats={statsCards} />
      )}

      {/* Pending Invoice Alert */}
      {pendingInvoice && (
        <Card className="border-amber-200 bg-amber-50">
          <CardContent className="flex items-center justify-between py-4">
            <div>
              <p className="font-semibold text-amber-900">
                Invoice Pending Payment
              </p>
              <p className="text-sm text-amber-700">
                {pendingInvoice.invoice_number} —{" "}
                {currency(pendingInvoice.total_amount)} due{" "}
                {pendingInvoice.due_date ?? "soon"}
              </p>
            </div>
            <Button asChild size="sm">
              <a href="/customer/payments">Pay Now</a>
            </Button>
          </CardContent>
        </Card>
      )}

      {/* Recent Order + Timeline */}
      <div className="grid gap-6 xl:grid-cols-[1.2fr_0.8fr]">
        <Card>
          <CardHeader>
            <CardTitle className="text-lg">Latest Order</CardTitle>
            <CardDescription>
              {ordersLoading
                ? "Loading..."
                : orders.length === 0
                  ? "No orders yet"
                  : `${orders.length} order${orders.length > 1 ? "s" : ""} total`}
            </CardDescription>
          </CardHeader>
          <CardContent>
            {ordersLoading ? (
              <div className="space-y-2">
                <div className="h-5 w-40 animate-pulse rounded bg-slate-200" />
                <div className="h-4 w-64 animate-pulse rounded bg-slate-200" />
                <div className="h-4 w-32 animate-pulse rounded bg-slate-100" />
              </div>
            ) : ordersError ? (
              <div className="flex items-center justify-between">
                <p className="text-sm text-red-600">Failed to load orders.</p>
                <Button size="sm" variant="outline" onClick={() => refetchOrders()}>
                  Retry
                </Button>
              </div>
            ) : latestOrder ? (
              <div className="space-y-3">
                <div className="flex items-start justify-between">
                  <div>
                    <p className="text-sm text-slate-500">
                      {latestOrder.order_number}
                    </p>
                    <p className="text-xl font-semibold text-slate-950">
                      {latestOrder.title}
                    </p>
                  </div>
                  <Badge tone={statusTone(latestOrder.status)}>
                    {latestOrder.status}
                  </Badge>
                </div>
                <p className="text-sm leading-6 text-slate-600">
                  {latestOrder.description}
                </p>
                <div className="flex flex-wrap gap-3 text-sm text-slate-500">
                  <span>{latestOrder.service_name}</span>
                  <span>•</span>
                  <span>{latestOrder.package_name}</span>
                  <span>•</span>
                  <span>{currency(latestOrder.total_price)}</span>
                </div>
                <Button
                  asChild
                  variant="outline"
                  size="sm"
                  className="mt-2"
                >
                  <a
                    href={`/customer/orders/${latestOrder.id}`}
                    onClick={() => setCurrentOrder(latestOrder)}
                  >
                    View Details
                  </a>
                </Button>
              </div>
            ) : (
              <EmptyState
                title="No orders yet"
                description="Create your first order to get started."
              />
            )}
          </CardContent>
        </Card>

        {/* Project Progress Timeline */}
        {latestOrder ? (
          <Card>
            <CardHeader>
              <CardTitle className="text-lg">Project Progress</CardTitle>
              <CardDescription>
                Milestones for {latestOrder.order_number}
              </CardDescription>
            </CardHeader>
            <CardContent>
              {progressLoading ? (
                <div className="space-y-3">
                  {[1, 2, 3].map((n) => (
                    <div key={n} className="flex gap-3">
                      <div className="h-4 w-4 animate-pulse rounded-full bg-slate-200" />
                      <div className="flex-1">
                        <div className="h-4 w-32 animate-pulse rounded bg-slate-200" />
                        <div className="mt-1 h-3 w-48 animate-pulse rounded bg-slate-100" />
                      </div>
                    </div>
                  ))}
                </div>
              ) : progress.length > 0 ? (
                <div className="space-y-4">
                  {progress.map((p, i) => (
                    <div key={p.id} className="flex gap-3">
                      <div className="flex flex-col items-center">
                        <div className="flex h-7 w-7 items-center justify-center rounded-full bg-blue-600 text-xs font-bold text-white">
                          {p.progress_percentage}%
                        </div>
                        {i < progress.length - 1 && (
                          <div className="mt-1 h-full w-px flex-1 bg-blue-100" />
                        )}
                      </div>
                      <div className="pb-4">
                        <p className="font-semibold text-slate-950">{p.title}</p>
                        <p className="text-sm text-slate-600">{p.description}</p>
                      </div>
                    </div>
                  ))}
                </div>
              ) : (
                <Timeline
                  currentIndex={statusIndex(latestOrder.status)}
                  steps={[
                    { title: "Requirement", description: "Brief submitted" },
                    { title: "Invoice", description: "Admin invoices" },
                    { title: "Payment", description: "Payment proof" },
                    { title: "In Progress", description: "Work started" },
                    { title: "Revision", description: "Feedback cycle" },
                    { title: "Completed", description: "Delivered" },
                  ]}
                />
              )}
            </CardContent>
          </Card>
        ) : (
          <Timeline
            currentIndex={0}
            steps={[
              { title: "Requirement", description: "Submit project brief" },
              { title: "Invoice", description: "Admin creates invoice" },
              { title: "Payment", description: "Upload payment proof" },
              { title: "In Progress", description: "Work starts" },
              { title: "Revision", description: "Feedback loop" },
              { title: "Completed", description: "Final delivery" },
            ]}
          />
        )}
      </div>

      {/* Quick Actions */}
      <Card>
        <CardHeader>
          <CardTitle className="text-lg">Quick Actions</CardTitle>
          <CardDescription>Jump into the main parts of your workflow.</CardDescription>
        </CardHeader>
        <CardContent className="flex flex-wrap gap-3">
          {[
            ["/customer/orders", "New Order"],
            ["/customer/invoices", "View Invoices"],
            ["/customer/payments", "Payments"],
            ["/customer/projects", "Project Timeline"],
            ["/customer/files", "Upload Files"],
          ].map(([href, label]) => (
            <Button key={href} asChild variant="outline">
              <a href={href}>{label}</a>
            </Button>
          ))}
        </CardContent>
      </Card>

      {/* Recent Activity */}
      <Card>
        <CardHeader>
          <CardTitle className="text-lg">Recent Activity</CardTitle>
          <CardDescription>Your latest actions on the platform.</CardDescription>
        </CardHeader>
        <CardContent>
          {activitiesLoading ? (
            <div className="space-y-3">
              {[1, 2, 3].map((n) => (
                <div
                  key={n}
                  className="h-10 w-full animate-pulse rounded-lg bg-slate-100"
                />
              ))}
            </div>
          ) : activities.length === 0 ? (
            <EmptyState
              title="No recent activity"
              description="Your actions will appear here as you use the platform."
            />
          ) : (
            <div className="divide-y divide-slate-100">
              {activities.slice(0, 10).map((log) => (
                <div
                  key={log.id}
                  className="flex items-start justify-between py-3"
                >
                  <div>
                    <p className="text-sm font-medium text-slate-900">
                      <span className="mr-2 inline-flex rounded bg-slate-100 px-1.5 py-0.5 text-xs text-slate-600">
                        {log.module}
                      </span>
                      {log.action}
                    </p>
                    {log.description && (
                      <p className="text-xs text-slate-500">{log.description}</p>
                    )}
                  </div>
                  <time className="ml-4 shrink-0 text-xs text-slate-400">
                    {new Date(log.created_at).toLocaleDateString("id-ID")}
                  </time>
                </div>
              ))}
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  );
}

export function CustomerOrdersView() {
  const router = useRouter();
  const setCurrentOrder = useProjectStore((state) => state.setCurrentOrder);
  const { data: services, isLoading: servicesLoading } = useQuery({
    queryKey: ["services"],
    queryFn: catalogApi.services,
  });
  const { data: packages, isLoading: packagesLoading } = useQuery({
    queryKey: ["packages"],
    queryFn: catalogApi.packages,
  });

  const createOrder = useMutation({
    mutationFn: orderApi.create,
    onSuccess: (order) => {
      setCurrentOrder(order);
      router.push(`/customer/orders/${order.id}`);
    },
  });

  const serviceList = services ?? [];
  const packageList = packages ?? [];

  return (
    <div className="space-y-8">
      <PageHeader
        eyebrow="Create Order"
        title="Order a new project"
        description="Choose a service and package, then submit the project brief for the existing backend flow."
      />
      <div className="grid gap-6 xl:grid-cols-[1fr_0.4fr]">
        <Card>
          <CardHeader>
            <CardTitle className="text-lg">Project Form</CardTitle>
            <CardDescription>Step 1 to 5 of the MVP flow.</CardDescription>
          </CardHeader>
          <CardContent>
            {servicesLoading || packagesLoading ? (
              <div className="flex items-center gap-2 text-sm text-slate-500">
                <Loader2 className="h-4 w-4 animate-spin" /> Loading services
                and packages...
              </div>
            ) : (
              <OrderForm
                services={serviceList.map((item) => ({
                  id: item.id,
                  name: item.name,
                  category_name: item.category_name,
                }))}
                packages={packageList.map((item) => ({
                  id: item.id,
                  name: item.name,
                  service_id: item.service_id,
                  price: item.price,
                  revision_count: item.revision_count,
                  delivery_days: item.delivery_days,
                }))}
                onSubmit={(values) => createOrder.mutate(values)}
              />
            )}
          </CardContent>
        </Card>
        <Card>
          <CardHeader>
            <CardTitle className="text-lg">Selection Summary</CardTitle>
            <CardDescription>
              Service and package data comes from the backend.
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            {(services ?? []).slice(0, 4).map((service) => (
              <div
                key={service.id}
                className="rounded-3xl border border-slate-200 p-4"
              >
                <p className="font-semibold text-slate-950">{service.name}</p>
                <p className="text-sm text-slate-500">
                  {service.category_name}
                </p>
              </div>
            ))}
          </CardContent>
        </Card>
      </div>
    </div>
  );
}

export function CustomerOrderDetailView({ orderId }: { orderId: string }) {
  const currentOrder = useProjectStore((state) => state.currentOrder);
  const requirements = useProjectStore((state) => state.requirements);
  const addRequirement = useProjectStore((state) => state.addRequirement);
  const currentInvoice = useProjectStore((state) => state.currentInvoice);
  const currentPayment = useProjectStore((state) => state.currentPayment);

  const saveRequirement = useMutation({
    mutationFn: (values: { question: string; answer: string }) =>
      requirementApi.create(orderId, values),
    onSuccess: (requirement) => addRequirement(requirement),
  });

  const matchesCurrent = currentOrder?.id === orderId;

  if (!matchesCurrent && !currentOrder) {
    return (
      <EmptyState
        title="Order details unavailable"
        description="Create a new order in this session to populate the detail view."
      />
    );
  }

  if (!matchesCurrent) {
    return (
      <EmptyState
        title="Order not loaded locally"
        description="Open the order from the same session to view the stored API response."
      />
    );
  }

  return (
    <div className="space-y-8">
      <PageHeader
        eyebrow="Order Detail"
        title={currentOrder.title}
        description={`Order ${currentOrder.order_number} is tracked here with the rest of the project flow.`}
      />
      <div className="grid gap-6 xl:grid-cols-2">
        <Card>
          <CardHeader>
            <CardTitle className="text-lg">Order Information</CardTitle>
          </CardHeader>
          <CardContent className="space-y-3 text-sm text-slate-600">
            <p>
              <span className="font-medium text-slate-900">Service:</span>{" "}
              {currentOrder.service_name}
            </p>
            <p>
              <span className="font-medium text-slate-900">Package:</span>{" "}
              {currentOrder.package_name}
            </p>
            <p>
              <span className="font-medium text-slate-900">Budget:</span>{" "}
              {currency(currentOrder.total_price)}
            </p>
            <p>
              <span className="font-medium text-slate-900">Status:</span>{" "}
              <Badge tone={statusTone(currentOrder.status)}>
                {currentOrder.status}
              </Badge>
            </p>
            <p className="leading-7">{currentOrder.description}</p>
          </CardContent>
        </Card>
        <Timeline
          steps={[
            { title: "Requirement", description: "Requirement submission" },
            { title: "Invoice", description: "Admin invoice created" },
            { title: "Payment", description: "Customer payment proof" },
            { title: "In Progress", description: "Work in delivery" },
            { title: "Revision", description: "Feedback loop" },
            { title: "Completed", description: "Final closeout" },
          ]}
          currentIndex={statusIndex(currentOrder.status)}
        />
      </div>

      <Card>
        <CardHeader>
          <CardTitle className="text-lg">Requirements</CardTitle>
          <CardDescription>
            Body uses question and answer exactly as the backend expects.
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-6">
          <RequirementForm
            onSubmit={(values) => saveRequirement.mutate(values)}
          />
          <div className="grid gap-3 md:grid-cols-2">
            {requirements.map((requirement) => (
              <div
                key={requirement.id}
                className="rounded-3xl border border-slate-200 p-4"
              >
                <p className="text-sm font-semibold text-slate-950">
                  {requirement.question}
                </p>
                <p className="mt-2 text-sm leading-6 text-slate-600">
                  {requirement.answer}
                </p>
              </div>
            ))}
          </div>
        </CardContent>
      </Card>

      <div className="grid gap-6 xl:grid-cols-2">
        <Card>
          <CardHeader>
            <CardTitle className="text-lg">Invoice</CardTitle>
          </CardHeader>
          <CardContent className="space-y-3 text-sm text-slate-600">
            {currentInvoice ? (
              <>
                <p>
                  <span className="font-medium text-slate-900">Invoice:</span>{" "}
                  {currentInvoice.invoice_number}
                </p>
                <p>
                  <span className="font-medium text-slate-900">Status:</span>{" "}
                  <Badge tone={statusTone(currentInvoice.status)}>
                    {currentInvoice.status}
                  </Badge>
                </p>
                <p>
                  <span className="font-medium text-slate-900">Due Date:</span>{" "}
                  {currentInvoice.due_date ?? "-"}
                </p>
                <p>
                  <span className="font-medium text-slate-900">Total:</span>{" "}
                  {currency(currentInvoice.total_amount)}
                </p>
              </>
            ) : (
              <EmptyState
                title="Invoice not available yet"
                description="The admin creates the invoice after the requirement step."
              />
            )}
          </CardContent>
        </Card>
        <Card>
          <CardHeader>
            <CardTitle className="text-lg">Payment</CardTitle>
          </CardHeader>
          <CardContent className="space-y-3 text-sm text-slate-600">
            {currentPayment ? (
              <>
                <p>
                  <span className="font-medium text-slate-900">Status:</span>{" "}
                  <Badge tone={statusTone(currentPayment.payment_status)}>
                    {currentPayment.payment_status}
                  </Badge>
                </p>
                <p>
                  <span className="font-medium text-slate-900">Method:</span>{" "}
                  {currentPayment.payment_method ?? "-"}
                </p>
                <p>
                  <span className="font-medium text-slate-900">Proof:</span>{" "}
                  {currentPayment.payment_proof_url ?? "-"}
                </p>
              </>
            ) : (
              <EmptyState
                title="Payment not submitted yet"
                description="After invoice creation, continue to the payment page."
              />
            )}
          </CardContent>
        </Card>
      </div>
    </div>
  );
}

export function CustomerInvoicesView() {
  const { data: invoices = [] } = useQuery({
    queryKey: ["my-invoices"],
    queryFn: invoiceApi.myInvoices,
  });

  if (invoices.length === 0) {
    return (
      <div className="space-y-8">
        <PageHeader
          eyebrow="Invoices"
          title="Your Invoices"
          description="Track all your invoices and their payment status."
        />
        <EmptyState
          title="No invoices yet"
          description="Invoices will appear after your orders are approved and billed."
        />
      </div>
    );
  }

  return (
    <div className="space-y-8">
      <PageHeader
        eyebrow="Invoices"
        title="Your Invoices"
        description="Track billing status for all your projects."
      />
      <div className="space-y-4">
        {invoices.map((invoice) => (
          <Card key={invoice.id}>
            <CardHeader>
              <div className="flex items-center justify-between">
                <div>
                  <CardTitle className="text-lg">
                    {invoice.invoice_number}
                  </CardTitle>
                  <CardDescription>Order {invoice.order_number}</CardDescription>
                </div>
                <Badge tone={statusTone(invoice.status)}>
                  {invoice.status}
                </Badge>
              </div>
            </CardHeader>
            <CardContent className="space-y-3 text-sm text-slate-600">
              <div className="grid gap-4 sm:grid-cols-2">
                <div>
                  <p className="font-medium text-slate-900">Subtotal:</p>
                  <p>{currency(invoice.subtotal)}</p>
                </div>
                <div>
                  <p className="font-medium text-slate-900">Total:</p>
                  <p>{currency(invoice.total_amount)}</p>
                </div>
                <div>
                  <p className="font-medium text-slate-900">Due Date:</p>
                  <p>{invoice.due_date ?? "-"}</p>
                </div>
                <div>
                  <p className="font-medium text-slate-900">Status:</p>
                  <p>{invoice.status}</p>
                </div>
              </div>
              {invoice.status === "UNPAID" && (
                <Button asChild>
                  <a href="/customer/payments">Pay Now</a>
                </Button>
              )}
            </CardContent>
          </Card>
        ))}
      </div>
    </div>
  );
}

export function CustomerInvoiceDetailView({
  invoiceId,
}: {
  invoiceId: string;
}) {
  const invoice = useProjectStore((state) => state.currentInvoice);
  if (!invoice || invoice.id !== invoiceId) {
    return (
      <EmptyState
        title="Invoice not loaded locally"
        description="Open the invoice from the same browser session to review it here."
      />
    );
  }
  return <CustomerInvoicesView />;
}

export function CustomerPaymentsView() {
  const user = useAuthStore((state) => state.user);
  const order = useProjectStore((state) => state.currentOrder);
  const invoice = useProjectStore((state) => state.currentInvoice);
  const setCurrentPayment = useProjectStore((state) => state.setCurrentPayment);
  const [paymentMethod, setPaymentMethod] = useState("");
  const [file, setFile] = useState<File | null>(null);
  const [message, setMessage] = useState<string | null>(null);

  const { data: payments = [], refetch: refetchPayments } = useQuery({
    queryKey: ["my-payments"],
    queryFn: paymentApi.myPayments,
  });

  const { data: invoices = [] } = useQuery({
    queryKey: ["my-invoices"],
    queryFn: invoiceApi.myInvoices,
  });

  const submitPayment = useMutation({
    mutationFn: async () => {
      if (!order || !invoice || !user)
        throw new Error("Order, invoice, and user are required");
      let proofUrl = "";
      if (file) {
        const formData = new FormData();
        formData.append("order_id", order.id);
        formData.append("uploaded_by", user.id);
        formData.append("file_type", "payment-proof");
        formData.append("file", file);
        const uploaded = await fileApi.upload(formData);
        proofUrl = uploaded.file_url;
      }

      return paymentApi.create({
        invoice_id: invoice.id,
        amount: invoice.total_amount,
        payment_method: paymentMethod,
        payment_proof_url: proofUrl,
      });
    },
    onSuccess: (payment) => {
      setCurrentPayment(payment);
      setMessage("Payment submitted successfully.");
      setPaymentMethod("");
      setFile(null);
      refetchPayments();
    },
  });

  return (
    <div className="space-y-8">
      <PageHeader
        eyebrow="Payments"
        title="Manage Your Payments"
        description="Upload payment proof or view payment history."
      />

      {invoice && invoice.status === "UNPAID" && (
        <Card>
          <CardHeader>
            <CardTitle className="text-lg">Submit New Payment</CardTitle>
            <CardDescription>
              Invoice: {invoice.invoice_number} - {currency(invoice.total_amount)}
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="grid gap-4 lg:grid-cols-2">
              <div>
                <label className="mb-2 block text-sm font-medium text-slate-700">
                  Payment Method
                </label>
                <Input
                  value={paymentMethod}
                  onChange={(event) => setPaymentMethod(event.target.value)}
                  placeholder="Bank Transfer / QRIS"
                />
              </div>
              <div>
                <label className="mb-2 block text-sm font-medium text-slate-700">
                  Upload Proof
                </label>
                <Input
                  type="file"
                  onChange={(event) =>
                    setFile(event.target.files?.[0] ?? null)
                  }
                />
              </div>
            </div>
            <Button
              onClick={() => submitPayment.mutate()}
              disabled={submitPayment.isPending}
            >
              {submitPayment.isPending ? (
                <Loader2 className="mr-2 h-4 w-4 animate-spin" />
              ) : (
                <Upload className="mr-2 h-4 w-4" />
              )}
              Submit Payment
            </Button>
            {message ? (
              <p className="text-sm text-emerald-600">{message}</p>
            ) : null}
          </CardContent>
        </Card>
      )}

      <Card>
        <CardHeader>
          <CardTitle className="text-lg">Payment History</CardTitle>
          <CardDescription>
            All payments for your projects
          </CardDescription>
        </CardHeader>
        <CardContent>
          {payments.length === 0 ? (
            <EmptyState
              title="No payments yet"
              description="Payment records will appear here after submission."
            />
          ) : (
            <div className="space-y-4">
              {payments.map((payment) => (
                <div
                  key={payment.id}
                  className="flex items-center justify-between rounded-lg border border-slate-200 p-4"
                >
                  <div>
                    <p className="font-medium text-slate-900">
                      {currency(payment.amount)}
                    </p>
                    <p className="text-sm text-slate-500">
                      {payment.payment_method ?? "-"}
                    </p>
                  </div>
                  <Badge tone={statusTone(payment.payment_status)}>
                    {payment.payment_status}
                  </Badge>
                </div>
              ))}
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  );
}

export function CustomerProjectsView() {
  const order = useProjectStore((state) => state.currentOrder);
  const invoice = useProjectStore((state) => state.currentInvoice);
  const payment = useProjectStore((state) => state.currentPayment);
  const currentIndex = useMemo(
    () =>
      statusIndex(order?.status ?? invoice?.status ?? payment?.payment_status),
    [invoice?.status, order?.status, payment?.payment_status],
  );

  return (
    <div className="space-y-8">
      <PageHeader
        eyebrow="Projects"
        title="Project timeline"
        description="A vertical timeline for the active order lifecycle."
      />
      <Timeline
        currentIndex={currentIndex}
        steps={[
          { title: "Requirement", description: "Requirement captured" },
          { title: "Invoice", description: "Invoice generated" },
          { title: "Payment", description: "Payment submitted" },
          { title: "In Progress", description: "Work started" },
          { title: "Revision", description: "Feedback cycle" },
          { title: "Completed", description: "Project delivered" },
        ]}
      />
    </div>
  );
}

export function CustomerFilesView() {
  const user = useAuthStore((state) => state.user);
  const order = useProjectStore((state) => state.currentOrder);
  const [file, setFile] = useState<File | null>(null);
  const [fileType, setFileType] = useState("supporting-file");
  const [uploadedUrl, setUploadedUrl] = useState<string | null>(null);

  const upload = useMutation({
    mutationFn: async () => {
      if (!file || !order || !user)
        throw new Error("Order and file are required");
      const formData = new FormData();
      formData.append("order_id", order.id);
      formData.append("uploaded_by", user.id);
      formData.append("file_type", fileType);
      formData.append("file", file);
      return fileApi.upload(formData);
    },
    onSuccess: (result) => setUploadedUrl(result.file_url),
  });

  if (!order) {
    return (
      <EmptyState
        title="No order selected"
        description="Create an order first to upload supporting files."
      />
    );
  }

  return (
    <div className="space-y-8">
      <PageHeader
        eyebrow="Files"
        title="Upload project files"
        description="Send logos, references, and supporting documents using the existing upload endpoint."
      />
      <Card>
        <CardHeader>
          <CardTitle className="text-lg">File Upload</CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="grid gap-4 lg:grid-cols-2">
            <Input
              type="file"
              onChange={(event) => setFile(event.target.files?.[0] ?? null)}
            />
            <Input
              value={fileType}
              onChange={(event) => setFileType(event.target.value)}
              placeholder="file type"
            />
          </div>
          <Button onClick={() => upload.mutate()} disabled={upload.isPending}>
            {upload.isPending ? (
              <Loader2 className="mr-2 h-4 w-4 animate-spin" />
            ) : (
              <Upload className="mr-2 h-4 w-4" />
            )}
            Upload File
          </Button>
          {uploadedUrl ? (
            <p className="text-sm text-emerald-600">
              Uploaded to {uploadedUrl}
            </p>
          ) : null}
        </CardContent>
      </Card>
    </div>
  );
}

export function CustomerMessagesView() {
  return (
    <EmptyState
      title="Messages"
      description="A messaging module can be attached when backend endpoints are ready."
    />
  );
}

export function CustomerProfileView() {
  const user = useAuthStore((state) => state.user);
  return (
    <div className="space-y-8">
      <PageHeader
        eyebrow="Profile"
        title="Your profile"
        description="Account details are read from the authenticated session."
      />
      <Card>
        <CardHeader>
          <CardTitle className="text-lg">Account Info</CardTitle>
        </CardHeader>
        <CardContent className="space-y-3 text-sm text-slate-600">
          <p>
            <span className="font-medium text-slate-900">Name:</span>{" "}
            {user?.name ?? "-"}
          </p>
          <p>
            <span className="font-medium text-slate-900">Email:</span>{" "}
            {user?.email ?? "-"}
          </p>
          <p>
            <span className="font-medium text-slate-900">Role:</span>{" "}
            {user?.role ?? "-"}
          </p>
        </CardContent>
      </Card>
    </div>
  );
}

export function CustomerSettingsView() {
  return (
    <EmptyState
      title="Settings"
      description="Notification and preference controls can be added without changing the backend contract."
    />
  );
}
