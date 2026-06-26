export type ApiEnvelope<T> = {
  success?: boolean;
  message?: string;
  data: T;
  error?: string;
};

export type AuthUser = {
  id: string;
  name: string;
  email: string;
  phone?: string;
  role: string;
  status?: string;
};

export type AuthResponse = {
  token: string;
  user: AuthUser;
};

export type LoginPayload = {
  email: string;
  password: string;
};

export type RegisterPayload = {
  name: string;
  email: string;
  password: string;
  phone?: string;
};

export type ServiceItem = {
  id: string;
  category_id: string;
  category_name: string;
  name: string;
  slug: string;
  description: string;
  base_price: number;
  estimated_days: number;
  status: string;
};

export type PackageItem = {
  id: string;
  service_id: string;
  service_name: string;
  name: string;
  description: string;
  price: number;
  revision_count: number;
  delivery_days: number;
  features: string[];
  status: string;
};

export type OrderItem = {
  id: string;
  customer_id: string;
  customer_name: string;
  service_id: string;
  service_name: string;
  package_id: string;
  package_name: string;
  order_number: string;
  title: string;
  description: string;
  deadline?: string | null;
  total_price: number;
  status: string;
};

export type RequirementItem = {
  id: string;
  order_id: string;
  question: string;
  answer: string;
};

export type InvoiceItem = {
  id: string;
  order_id: string;
  order_number: string;
  invoice_number: string;
  subtotal: number;
  discount: number;
  tax: number;
  total_amount: number;
  status: string;
  due_date?: string | null;
};

export type PaymentItem = {
  id: string;
  invoice_id: string;
  amount: number;
  payment_method?: string | null;
  payment_status: string;
  payment_proof_url?: string | null;
  paid_at?: string | null;
  verified_by?: string | null;
};

export type DashboardStats = {
  total_customers: number;
  total_orders: number;
  total_services: number;
  total_packages: number;
  pending_orders: number;
  in_progress_orders: number;
  completed_orders: number;
  pending_payments: number;
  paid_payments: number;
  total_revenue: number;
};

export type ActivityLogItem = {
  id: string;
  user_id: string;
  module: string;
  action: string;
  description: string;
  ip_address?: string | null;
  created_at: string;
};