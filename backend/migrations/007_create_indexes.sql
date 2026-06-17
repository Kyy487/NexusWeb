CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_role_id ON users(role_id);
CREATE INDEX idx_users_status ON users(status);

CREATE INDEX idx_services_category_id ON services(category_id);
CREATE INDEX idx_services_slug ON services(slug);
CREATE INDEX idx_services_status ON services(status);

CREATE INDEX idx_service_packages_service_id ON service_packages(service_id);

CREATE INDEX idx_orders_customer_id ON service_orders(customer_id);
CREATE INDEX idx_orders_service_id ON service_orders(service_id);
CREATE INDEX idx_orders_status ON service_orders(status);
CREATE INDEX idx_orders_order_number ON service_orders(order_number);

CREATE INDEX idx_invoices_order_id ON invoices(order_id);
CREATE INDEX idx_invoices_status ON invoices(status);

CREATE INDEX idx_payments_invoice_id ON payments(invoice_id);
CREATE INDEX idx_payments_status ON payments(payment_status);

CREATE INDEX idx_messages_order_id ON messages(order_id);
CREATE INDEX idx_testimonials_service_id ON testimonials(service_id);
CREATE INDEX idx_activity_logs_user_id ON activity_logs(user_id);