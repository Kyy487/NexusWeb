"use client";

import { useMutation, useQuery } from "@tanstack/react-query";
import { Loader2 } from "lucide-react";
import { useState } from "react";

import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { DataTable } from "@/components/dashboard/data-table";
import { EmptyState } from "@/components/dashboard/empty-state";
import { PageHeader } from "@/components/dashboard/page-header";
import { StatsGrid } from "@/components/dashboard/stats-grid";
import {
  catalogApi,
  dashboardApi,
  invoiceApi,
  orderApi,
  paymentApi,
  userApi,
} from "@/lib/api";
import { useProjectStore } from "@/store/project-store";

function currency(value?: number | null) {
  return new Intl.NumberFormat("id-ID", {
    style: "currency",
    currency: "IDR",
    maximumFractionDigits: 0,
  }).format(value ?? 0);
}

function tone(status?: string) {
  if (!status) return "default";
  if (["PAID", "COMPLETED", "ACTIVE", "IN_PROGRESS"].includes(status))
    return "success";
  if (["PENDING", "UNPAID", "OVERDUE", "REVISION"].includes(status))
    return "warning";
  if (["CANCELLED", "FAILED", "EXPIRED"].includes(status)) return "danger";
  return "info";
}

export function AdminDashboardView() {
  const { data } = useQuery({
    queryKey: ["admin-stats"],
    queryFn: dashboardApi.stats,
  });
  const { data: activities } = useQuery({
    queryKey: ["admin-activities"],
    queryFn: dashboardApi.activities,
  });

  if (!data) {
    return (
      <EmptyState
        title="Loading dashboard"
        description="Fetching administrative metrics from the backend."
      />
    );
  }

  return (
    <div className="space-y-8">
      <PageHeader
        eyebrow="Admin Dashboard"
        title="Operations overview"
        description="Revenue, orders, customers, and recent activity in a premium command center."
      />
      <StatsGrid
        stats={[
          {
            label: "Revenue",
            value: currency(data.total_revenue),
            note: "Total revenue",
            tone: "success",
          },
          {
            label: "Orders",
            value: String(data.total_orders),
            note: `${data.pending_orders} pending`,
            tone: "info",
          },
          {
            label: "Customers",
            value: String(data.total_customers),
            note: "Registered users",
            tone: "info",
          },
          {
            label: "Projects",
            value: String(data.completed_orders),
            note: "Completed projects",
            tone: "success",
          },
        ]}
      />
      <div className="grid gap-6 xl:grid-cols-2">
        <Card>
          <CardHeader>
            <CardTitle className="text-lg">Recent Activities</CardTitle>
            <CardDescription>
              Audit trail captured by the backend.
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-3 text-sm text-slate-600">
            {(activities ?? []).slice(0, 8).map((item) => (
              <div
                key={item.id}
                className="rounded-3xl border border-slate-200 p-4"
              >
                <p className="font-semibold text-slate-950">
                  {item.module} - {item.action}
                </p>
                <p>{item.description}</p>
              </div>
            ))}
          </CardContent>
        </Card>
        <Card>
          <CardHeader>
            <CardTitle className="text-lg">Pending Payments</CardTitle>
          </CardHeader>
          <CardContent className="space-y-3 text-sm text-slate-600">
            <p>Pending: {data.pending_payments}</p>
            <p>Paid: {data.paid_payments}</p>
            <div className="h-3 rounded-full bg-slate-100">
              <div
                className="h-3 rounded-full bg-blue-600"
                style={{
                  width: `${Math.max(10, (data.paid_payments / Math.max(1, data.pending_payments + data.paid_payments)) * 100)}%`,
                }}
              />
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}

export function AdminOrdersView() {
  const { data, isLoading } = useQuery({
    queryKey: ["admin-orders"],
    queryFn: orderApi.list,
  });
  const updateStatus = useMutation({
    mutationFn: ({ id, status }: { id: string; status: string }) =>
      orderApi.updateStatus(id, status),
  });

  if (isLoading)
    return (
      <EmptyState
        title="Loading orders"
        description="Fetching order list from the backend."
      />
    );

  return (
    <div className="space-y-8">
      <PageHeader
        eyebrow="Orders"
        title="Manage customer orders"
        description="Browse the existing service orders and update their status."
      />
      <DataTable
        title="All Orders"
        description="Backend-provided order records."
        columns={["Order Number", "Title", "Status", "Total Price"]}
        rows={(data ?? []).map((order) => ({
          id: order.id,
          order_number: order.order_number,
          title: order.title,
          status: order.status,
          total_price: currency(order.total_price),
        }))}
      />
      <Card>
        <CardHeader>
          <CardTitle className="text-lg">Update Order Status</CardTitle>
          <CardDescription>
            Choose a status that already exists in the backend enum.
          </CardDescription>
        </CardHeader>
        <CardContent className="grid gap-4 lg:grid-cols-3">
          <Input placeholder="Order ID" id="order-id" />
          <Input placeholder="Status e.g. IN_PROGRESS" id="order-status" />
          <Button
            onClick={() => {
              const id =
                (document.getElementById("order-id") as HTMLInputElement | null)
                  ?.value ?? "";
              const status =
                (
                  document.getElementById(
                    "order-status",
                  ) as HTMLInputElement | null
                )?.value ?? "";
              updateStatus.mutate({ id, status });
            }}
          >
            Update Status
          </Button>
        </CardContent>
      </Card>
    </div>
  );
}

export function AdminInvoicesView() {
  const { data, isLoading } = useQuery({
    queryKey: ["admin-invoices"],
    queryFn: invoiceApi.list,
  });
  const [orderId, setOrderId] = useState("");
  const [discount, setDiscount] = useState("0");
  const [tax, setTax] = useState("0");
  const setCurrentInvoice = useProjectStore((state) => state.setCurrentInvoice);
  const createInvoice = useMutation({
    mutationFn: () =>
      invoiceApi.create({
        order_id: orderId,
        discount: Number(discount),
        tax: Number(tax),
      }),
    onSuccess: (invoice) => setCurrentInvoice(invoice),
  });

  if (isLoading)
    return (
      <EmptyState
        title="Loading invoices"
        description="Fetching invoice list from the backend."
      />
    );

  return (
    <div className="space-y-8">
      <PageHeader
        eyebrow="Invoices"
        title="Create and manage invoices"
        description="Create a new invoice for an existing order, then review all invoice records."
      />
      <Card>
        <CardHeader>
          <CardTitle className="text-lg">Create Invoice</CardTitle>
        </CardHeader>
        <CardContent className="grid gap-4 lg:grid-cols-4">
          <Input
            value={orderId}
            onChange={(event) => setOrderId(event.target.value)}
            placeholder="Order ID"
          />
          <Input
            value={discount}
            onChange={(event) => setDiscount(event.target.value)}
            placeholder="Discount"
          />
          <Input
            value={tax}
            onChange={(event) => setTax(event.target.value)}
            placeholder="Tax"
          />
          <Button
            onClick={() => createInvoice.mutate()}
            disabled={createInvoice.isPending}
          >
            {createInvoice.isPending ? (
              <Loader2 className="mr-2 h-4 w-4 animate-spin" />
            ) : null}
            Create Invoice
          </Button>
        </CardContent>
      </Card>
      <DataTable
        title="All Invoices"
        description="Backend-provided invoice records."
        columns={["Invoice Number", "Status", "Total Amount", "Due Date"]}
        rows={(data ?? []).map((invoice) => ({
          id: invoice.id,
          invoice_number: invoice.invoice_number,
          status: invoice.status,
          total_amount: currency(invoice.total_amount),
          due_date: invoice.due_date ?? "-",
        }))}
      />
    </div>
  );
}

export function AdminPaymentsView() {
  const { data, isLoading } = useQuery({
    queryKey: ["admin-payments"],
    queryFn: paymentApi.list,
  });
  const updateStatus = useMutation({
    mutationFn: ({
      id,
      status,
      verifiedBy,
    }: {
      id: string;
      status: string;
      verifiedBy?: string;
    }) => paymentApi.updateStatus(id, status, verifiedBy),
  });

  if (isLoading)
    return (
      <EmptyState
        title="Loading payments"
        description="Fetching payment records from the backend."
      />
    );

  return (
    <div className="space-y-8">
      <PageHeader
        eyebrow="Payments"
        title="Verify payment records"
        description="Approve or review submitted payments in the backend-driven workflow."
      />
      <DataTable
        title="All Payments"
        description="Payment records fetched from the API."
        columns={["Invoice Id", "Amount", "Status", "Proof"]}
        rows={(data ?? []).map((payment) => ({
          id: payment.id,
          invoice_id: payment.invoice_id,
          amount: currency(payment.amount),
          status: payment.payment_status,
          proof: payment.payment_proof_url ?? "-",
        }))}
      />
      <Card>
        <CardHeader>
          <CardTitle className="text-lg">Verify Payment</CardTitle>
        </CardHeader>
        <CardContent className="grid gap-4 lg:grid-cols-4">
          <Input placeholder="Payment ID" id="payment-id" />
          <Input placeholder="Payment Status e.g. PAID" id="payment-status" />
          <Input placeholder="Verified By user ID" id="verified-by" />
          <Button
            onClick={() => {
              const id =
                (
                  document.getElementById(
                    "payment-id",
                  ) as HTMLInputElement | null
                )?.value ?? "";
              const status =
                (
                  document.getElementById(
                    "payment-status",
                  ) as HTMLInputElement | null
                )?.value ?? "";
              const verifiedBy =
                (
                  document.getElementById(
                    "verified-by",
                  ) as HTMLInputElement | null
                )?.value ?? "";
              updateStatus.mutate({ id, status, verifiedBy });
            }}
          >
            Verify Payment
          </Button>
        </CardContent>
      </Card>
    </div>
  );
}

export function AdminCustomersView() {
  const { data, isLoading } = useQuery({
    queryKey: ["admin-customers"],
    queryFn: userApi.list,
  });
  if (isLoading)
    return (
      <EmptyState
        title="Loading customers"
        description="Fetching customer records from the backend."
      />
    );
  return (
    <div className="space-y-8">
      <PageHeader
        eyebrow="Customers"
        title="Customer directory"
        description="Track users registered in the system."
      />
      <DataTable
        title="Customers"
        description="Users fetched from /users."
        columns={["Name", "Email", "Role", "Status"]}
        rows={(data ?? []).map((user) => ({
          id: user.id,
          name: user.name,
          email: user.email,
          role: user.role,
          status: user.status ?? "-",
        }))}
      />
    </div>
  );
}

export function AdminServicesView() {
  const { data, isLoading } = useQuery({
    queryKey: ["admin-services"],
    queryFn: catalogApi.services,
  });
  if (isLoading)
    return (
      <EmptyState
        title="Loading services"
        description="Fetching service catalog from the backend."
      />
    );
  return (
    <DataTable
      title="Services"
      description="Service catalog from the backend."
      columns={["Name", "Category", "Price", "Status"]}
      rows={(data ?? []).map((item) => ({
        id: item.id,
        name: item.name,
        category: item.category_name,
        price: currency(item.base_price),
        status: item.status,
      }))}
    />
  );
}

export function AdminPackagesView() {
  const { data, isLoading } = useQuery({
    queryKey: ["admin-packages"],
    queryFn: catalogApi.packages,
  });
  if (isLoading)
    return (
      <EmptyState
        title="Loading packages"
        description="Fetching package catalog from the backend."
      />
    );
  return (
    <DataTable
      title="Packages"
      description="Package catalog from the backend."
      columns={["Name", "Service", "Price", "Status"]}
      rows={(data ?? []).map((item) => ({
        id: item.id,
        name: item.name,
        service: item.service_name,
        price: currency(item.price),
        status: item.status,
      }))}
    />
  );
}

export function AdminReportsView() {
  return (
    <EmptyState
      title="Reports"
      description="Report summaries and charts can be layered on top of the existing statistics endpoint."
    />
  );
}

export function AdminSettingsView() {
  return (
    <EmptyState
      title="Settings"
      description="Administrative settings can be added without changing the existing data contract."
    />
  );
}
