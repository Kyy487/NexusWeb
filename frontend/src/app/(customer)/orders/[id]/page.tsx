import { CustomerOrderDetailView } from "@/components/customer/customer-pages";

export default async function CustomerOrderDetailPage({ params }: Readonly<{ params: Promise<{ id: string }> }>) {
  const { id } = await params;
  return <CustomerOrderDetailView orderId={id} />;
}