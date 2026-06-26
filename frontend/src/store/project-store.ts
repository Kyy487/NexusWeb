import { create } from "zustand";
import { createJSONStorage, persist } from "zustand/middleware";

import type { InvoiceItem, OrderItem, PaymentItem, RequirementItem } from "@/types/api";

type ProjectFlowState = {
  currentOrder: OrderItem | null;
  currentInvoice: InvoiceItem | null;
  currentPayment: PaymentItem | null;
  requirements: RequirementItem[];
  setCurrentOrder: (order: OrderItem | null) => void;
  setCurrentInvoice: (invoice: InvoiceItem | null) => void;
  setCurrentPayment: (payment: PaymentItem | null) => void;
  addRequirement: (requirement: RequirementItem) => void;
  reset: () => void;
};

export const useProjectStore = create<ProjectFlowState>()(
  persist(
    (set) => ({
      currentOrder: null,
      currentInvoice: null,
      currentPayment: null,
      requirements: [],
      setCurrentOrder: (currentOrder) => set({ currentOrder }),
      setCurrentInvoice: (currentInvoice) => set({ currentInvoice }),
      setCurrentPayment: (currentPayment) => set({ currentPayment }),
      addRequirement: (requirement) =>
        set((state) => ({ requirements: [...state.requirements, requirement] })),
      reset: () => set({ currentOrder: null, currentInvoice: null, currentPayment: null, requirements: [] }),
    }),
    {
      name: "nexusweb-project-flow",
      storage: createJSONStorage(() => localStorage),
    }
  )
);