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
  fileApi,
  invoiceApi,
  orderApi,
  paymentApi,
  requirementApi,
} from "@/lib/api";
import { useAuthStore } from "@/store/auth-store";
import { useProjectStore } from "@/store/project-store";
import type { InvoiceItem, OrderItem, PaymentItem } from "@/types/api";

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
  const currentOrder = useProjectStore((state) => state.currentOrder);
  const currentInvoice = useProjectStore((state) => state.currentInvoice);
  const currentPayment = useProjectStore((state) => state.currentPayment);
  const requirements = useProjectStore((state) => state.requirements);

  const stats = [
    {
      label: "Total Orders",
      value: currentOrder ? "1" : "0",
      note: currentOrder?.status ?? "No orders yet",
      tone: statusTone(currentOrder?.status),
    },
    {
      label: "Active Projects",
      value:
        currentOrder &&
        !["COMPLETED", "CANCELLED"].includes(currentOrder.status)
          ? "1"
          : "0",
      note: currentOrder?.status ?? "Idle",
      tone: "info",
    },
    {
      label: "Pending Payments",
      value: currentPayment ? "0" : currentInvoice ? "1" : "0",
      note: currentInvoice?.status ?? "Awaiting invoice",
      tone: statusTone(currentInvoice?.status),
    },
    {
      label: "Completed Projects",
      value: currentOrder?.status === "COMPLETED" ? "1" : "0",
      note: currentPayment?.payment_status ?? "In progress",
      tone: statusTone(currentPayment?.payment_status),
    },
  ];

  return (
    <div className="space-y-8">
      <PageHeader
        eyebrow="Customer Dashboard"
        title={`Welcome${user?.name ? `, ${user.name}` : ""}`}
        description="Track your project lifecycle from order submission to completion in one premium workspace."
      />
      <StatsGrid stats={stats} />
      <div className="grid gap-6 xl:grid-cols-[1.2fr_0.8fr]">
        <Card>
          <CardHeader>
            <CardTitle className="text-lg">Recent Order</CardTitle>
            <CardDescription>
              Stored from the latest successful project submission.
            </CardDescription>
          </CardHeader>
          <CardContent>
            {currentOrder ? (
              <div className="space-y-3">
                <div className="flex items-center justify-between">
                  <div>
                    <p className="text-sm text-slate-500">
                      {currentOrder.order_number}
                    </p>
                    <p className="text-xl font-semibold text-slate-950">
                      {currentOrder.title}
                    </p>
                  </div>
                  <Badge tone={statusTone(currentOrder.status)}>
                    {" "}
                    {currentOrder.status}
                  </Badge>
                </div>
                <p className="text-sm leading-6 text-slate-600">
                  {currentOrder.description}
                </p>
                <div className="flex flex-wrap gap-3 text-sm text-slate-500">
                  <span>{currentOrder.service_name}</span>
                  <span>•</span>
                  <span>{currentOrder.package_name}</span>
                  <span>•</span>
                  <span>{currency(currentOrder.total_price)}</span>
                </div>
              </div>
            ) : (
              <EmptyState
                title="No project yet"
                description="Create your first order to populate the dashboard."
              />
            )}
          </CardContent>
        </Card>
        <Timeline
          steps={[
            {
              title: "Requirement",
              description: "Project brief is submitted.",
            },
            { title: "Invoice", description: "Admin creates the invoice." },
            {
              title: "Payment",
              description: "Customer submits payment proof.",
            },
            {
              title: "In Progress",
              description: "Work starts after verification.",
            },
            { title: "Revision", description: "Feedback and iteration." },
            { title: "Completed", description: "Final delivery and closeout." },
          ]}
          currentIndex={statusIndex(
            currentOrder?.status ??
              currentInvoice?.status ??
              currentPayment?.payment_status,
          )}
        />
      </div>
      <Card>
        <CardHeader>
          <CardTitle className="text-lg">Quick Actions</CardTitle>
          <CardDescription>
            Jump into the main parts of the MVP flow.
          </CardDescription>
        </CardHeader>
        <CardContent className="flex flex-wrap gap-3">
          {[
            ["/customer/orders", "Create Order"],
            ["/customer/invoices", "View Invoice"],
            ["/customer/payments", "Create Payment"],
            ["/customer/files", "Upload Files"],
          ].map(([href, label]) => (
            <Button key={href} asChild variant="outline">
              <a href={href}>{label}</a>
            </Button>
          ))}
        </CardContent>
      </Card>
      <DataTable
        title="Latest Requirements"
        description="Persisted from the last successful requirement submission."
        columns={["Question", "Answer"]}
        rows={requirements.map((requirement) => ({
          id: requirement.id,
          Question: requirement.question,
          Answer: requirement.answer,
        }))}
      />
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
  const invoice = useProjectStore((state) => state.currentInvoice);
  return invoice ? (
    <div className="space-y-8">
      <PageHeader
        eyebrow="Invoices"
        title={invoice.invoice_number}
        description="Track billing status from the current session."
      />
      <Card>
        <CardHeader>
          <CardTitle className="text-lg">Invoice Overview</CardTitle>
        </CardHeader>
        <CardContent className="space-y-3 text-sm text-slate-600">
          <p>
            <span className="font-medium text-slate-900">Status:</span>{" "}
            <Badge tone={statusTone(invoice.status)}>{invoice.status}</Badge>
          </p>
          <p>
            <span className="font-medium text-slate-900">Due:</span>{" "}
            {invoice.due_date ?? "-"}
          </p>
          <p>
            <span className="font-medium text-slate-900">Subtotal:</span>{" "}
            {currency(invoice.subtotal)}
          </p>
          <p>
            <span className="font-medium text-slate-900">Total:</span>{" "}
            {currency(invoice.total_amount)}
          </p>
          <Button asChild>
            <a href="/customer/payments">Pay Invoice</a>
          </Button>
        </CardContent>
      </Card>
    </div>
  ) : (
    <EmptyState
      title="No invoice in this session"
      description="The invoice will appear after the admin creates it for the active order."
    />
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
    },
  });

  if (!order || !invoice) {
    return (
      <EmptyState
        title="Payment flow unavailable"
        description="Create the order and invoice in the current session first."
      />
    );
  }

  return (
    <div className="space-y-8">
      <PageHeader
        eyebrow="Payments"
        title="Submit payment proof"
        description="Upload proof and create the payment record using the existing backend endpoints."
      />
      <Card>
        <CardHeader>
          <CardTitle className="text-lg">Payment Status</CardTitle>
          <CardDescription>Invoice: {invoice.invoice_number}</CardDescription>
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
                onChange={(event) => setFile(event.target.files?.[0] ?? null)}
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
      <Card>
        <CardHeader>
          <CardTitle className="text-lg">Payment History</CardTitle>
        </CardHeader>
        <CardContent>
          <EmptyState
            title="No prior payments in this session"
            description="Submitted payment records will appear after the API call succeeds."
          />
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
