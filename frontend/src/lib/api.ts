import axios from "axios";

import { useAuthStore } from "@/store/auth-store";
import type {
  ActivityLogItem,
  ApiEnvelope,
  AuthResponse,
  CustomerDashboardStats,
  DashboardStats,
  FileItem,
  InvoiceItem,
  LoginPayload,
  OrderItem,
  PackageItem,
  PaymentItem,
  ProgressItem,
  RegisterPayload,
  RequirementItem,
  ServiceItem,
  UserItem,
} from "@/types/api";

const baseURL =
  process.env.NEXT_PUBLIC_API_URL ??
  "http://localhost:8080/api/v1";

export const api = axios.create({
  baseURL,
  headers: {
    "Content-Type": "application/json",
  },
});

api.interceptors.request.use((config) => {
  const token = useAuthStore.getState().token;

  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }

  return config;
});

api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error?.response?.status === 401) {
      useAuthStore.getState().clearAuth();
    }

    return Promise.reject(error);
  }
);

async function unwrap<T>(
  promise: Promise<{ data: ApiEnvelope<T> }>
) {
  const response = await promise;
  return response.data.data;
}

export const authApi = {
  login: (payload: LoginPayload) =>
    unwrap<AuthResponse>(
      api.post("/auth/login", payload)
    ),

  register: (payload: RegisterPayload) =>
    unwrap<AuthResponse>(
      api.post("/auth/register", payload)
    ),
};

export const dashboardApi = {
  stats: () =>
    unwrap<DashboardStats>(
      api.get("/dashboard/stats")
    ),

  customerStats: () =>
    unwrap<CustomerDashboardStats>(
      api.get("/dashboard/stats")
    ),

  activities: () =>
    unwrap<ActivityLogItem[]>(
      api.get("/activity-logs")
    ),
};

export const catalogApi = {
  services: () =>
    unwrap<ServiceItem[]>(
      api.get("/services")
    ),

  packages: () =>
    unwrap<PackageItem[]>(
      api.get("/packages")
    ),
};

export const userApi = {
  me: () =>
    unwrap<UserItem>(
      api.get("/users/me")
    ),

  list: () =>
    unwrap<UserItem[]>(
      api.get("/users")
    ),
};

export const orderApi = {
  list: () =>
    unwrap<OrderItem[]>(
      api.get("/orders")
    ),

  // Deprecated: use list() — backend now filters by role automatically
  myOrders: () =>
    unwrap<OrderItem[]>(
      api.get("/orders")
    ),

  byId: (id: string) =>
    unwrap<OrderItem>(
      api.get(`/orders/${id}`)
    ),

  create: (payload: Record<string, unknown>) =>
    unwrap<OrderItem>(
      api.post("/orders", payload)
    ),

  updateStatus: (
    id: string,
    status: string
  ) =>
    unwrap<OrderItem>(
      api.patch(`/orders/${id}/status`, {
        status,
      })
    ),
};

export const requirementApi = {
  byOrderId: (orderId: string) =>
    unwrap<RequirementItem[]>(
      api.get(`/order-requirements/order/${orderId}`)
    ),

  create: (
    orderId: string,
    payload: {
      question: string;
      answer: string;
    }
  ) =>
    unwrap<RequirementItem>(
      api.post(
        `/order-requirements/order/${orderId}`,
        payload
      )
    ),
};

export const invoiceApi = {
  list: () =>
    unwrap<InvoiceItem[]>(
      api.get("/invoices")
    ),

  // Deprecated: use list() — backend now filters by role automatically
  myInvoices: () =>
    unwrap<InvoiceItem[]>(
      api.get("/invoices")
    ),

  byOrderId: (orderId: string) =>
    unwrap<InvoiceItem>(
      api.get(`/invoices/order/${orderId}`)
    ),

  byId: (id: string) =>
    unwrap<InvoiceItem>(
      api.get(`/invoices/${id}`)
    ),

  create: (payload: Record<string, unknown>) =>
    unwrap<InvoiceItem>(
      api.post("/invoices", payload)
    ),

  updateStatus: (
    id: string,
    status: string
  ) =>
    unwrap<InvoiceItem>(
      api.patch(`/invoices/${id}/status`, {
        status,
      })
    ),
};

export const paymentApi = {
  list: () =>
    unwrap<PaymentItem[]>(
      api.get("/payments")
    ),

  // Deprecated: use list() — backend now filters by role automatically
  myPayments: () =>
    unwrap<PaymentItem[]>(
      api.get("/payments")
    ),

  byInvoiceId: (invoiceId: string) =>
    unwrap<PaymentItem[]>(
      api.get(`/payments/invoice/${invoiceId}`)
    ),

  byId: (id: string) =>
    unwrap<PaymentItem>(
      api.get(`/payments/${id}`)
    ),

  create: (payload: Record<string, unknown>) =>
    unwrap<PaymentItem>(
      api.post("/payments", payload)
    ),

  updateStatus: (
    id: string,
    paymentStatus: string,
    verifiedBy?: string
  ) =>
    unwrap<PaymentItem>(
      api.patch(`/payments/${id}/status`, {
        payment_status: paymentStatus,
        verified_by: verifiedBy ?? "",
      })
    ),

  whatsapp: (id: string) =>
    unwrap<{
      invoice_id: string;
      amount: number;
      whatsapp_url: string;
    }>(
      api.get(`/payments/${id}/whatsapp`)
    ),
};

export const fileApi = {
  byOrderId: (orderId: string) =>
    unwrap<FileItem[]>(
      api.get(`/order-files/order/${orderId}`)
    ),

  upload: async (payload: FormData) => {
    const response = await api.post(
      "/order-files/upload",
      payload,
      {
        headers: {
          "Content-Type":
            "multipart/form-data",
        },
      }
    );

    return response.data.data;
  },
};

export const meApi = {
  me: () =>
    unwrap<UserItem>(
      api.get("/users/me")
    ),
};

export const progressApi = {
  byOrderId: (orderId: string) =>
    unwrap<ProgressItem[]>(
      api.get(`/order-progress/order/${orderId}`)
    ),
};