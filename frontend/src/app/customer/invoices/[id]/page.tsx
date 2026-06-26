import { CustomerInvoiceDetailView } from "@/components/customer/customer-pages";

export default async function CustomerInvoiceDetailPage({ params }: Readonly<{ params: Promise<{ id: string }> }>) {
  const { id } = await params;
  return <CustomerInvoiceDetailView invoiceId={id} />;
}