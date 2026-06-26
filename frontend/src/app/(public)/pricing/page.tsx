import { ContentPage } from "@/components/marketing/content-page";

export default function PricingPage() {
  return (
    <ContentPage
      title="Transparent pricing that builds trust"
      description="Pricing is presented simply so customers can understand scope before starting a project."
      points={[
        { title: "Transparent scope", description: "Clear assumptions and predictable project steps." },
        { title: "Flexible tiers", description: "Useful for different business sizes and timelines." },
        { title: "Secure payment flow", description: "Invoice, payment, and verification are handled in-platform." },
      ]}
    />
  );
}