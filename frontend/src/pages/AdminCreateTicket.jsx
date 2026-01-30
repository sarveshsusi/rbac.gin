import { useEffect, useState } from "react";
import { adminCreateTicket } from "../api/ticket.api";
import api from "../api/axios";

export default function AdminCreateTicket() {
    const [customers, setCustomers] = useState([]);
    const [customerProducts, setCustomerProducts] = useState([]);

    const [selectedCustomer, setSelectedCustomer] = useState("");
    const [selectedProduct, setSelectedProduct] = useState("");
    const [title, setTitle] = useState("");
    const [description, setDescription] = useState("");
    const [priority, setPriority] = useState("Standard");
    const [status, setStatus] = useState("Open");

    // Load Customers
    useEffect(() => {
        api.get("/admin/users?role=customer")
            .then(res => {
                // Handle both flat array (new API) and paginated object (old/defensive)
                const data = res.data;
                if (Array.isArray(data)) {
                    setCustomers(data);
                } else if (data.users && Array.isArray(data.users)) {
                    setCustomers(data.users);
                } else {
                    setCustomers([]);
                }
            })
            .catch(err => console.error(err));
    }, []);

    // Load Products when Customer selected
    useEffect(() => {
        if (!selectedCustomer) {
            setCustomerProducts([]);
            return;
        }
        api.get(`/admin/customers/${selectedCustomer}/products`)
            .then(res => {
                const data = res.data;
                if (Array.isArray(data)) {
                    setCustomerProducts(data);
                } else {
                    setCustomerProducts([]);
                }
            })
            .catch(() => setCustomerProducts([]));
    }, [selectedCustomer]);

    const submit = async (e) => {
        e.preventDefault();
        if (!selectedCustomer || !selectedProduct) return alert("Select customer and product");

        try {
            await adminCreateTicket({
                customer_id: selectedCustomer,
                product_id: selectedProduct,
                title,
                description,
                priority,
                status // Admin can set status or default is Open
            });
            alert("Ticket created for customer");
            setTitle("");
            setDescription("");
        } catch (err) {
            alert("Failed: " + (err.response?.data?.error || err.message));
        }
    };

    return (
        <div className="max-w-2xl bg-white p-6 rounded-xl shadow">
            <h1 className="text-xl font-bold mb-6">Create Ticket on Behalf</h1>

            <form onSubmit={submit} className="space-y-4">
                <div>
                    <label className="block text-sm font-medium mb-1">Customer</label>
                    <select
                        className="w-full border rounded p-2"
                        value={selectedCustomer}
                        onChange={e => setSelectedCustomer(e.target.value)}
                        required
                    >
                        <option value="">Select Customer</option>
                        {customers.map(c => (
                            <option key={c.id} value={c.id}>{c.name} ({c.email})</option>
                        ))}
                    </select>
                </div>

                <div>
                    <label className="block text-sm font-medium mb-1">Product</label>
                    <select
                        className="w-full border rounded p-2"
                        value={selectedProduct}
                        onChange={e => setSelectedProduct(e.target.value)}
                        required
                        disabled={!selectedCustomer}
                    >
                        <option value="">Select Product</option>
                        {Array.isArray(customerProducts) && customerProducts.map(p => (
                            <option key={p.id} value={p.product_id}>{p.product_name}</option>
                        ))}
                    </select>
                </div>

                <div>
                    <label className="block text-sm font-medium mb-1">Title</label>
                    <input
                        className="w-full border rounded p-2"
                        value={title}
                        onChange={e => setTitle(e.target.value)}
                        required
                    />
                </div>

                <div>
                    <label className="block text-sm font-medium mb-1">Description</label>
                    <textarea
                        className="w-full border rounded p-2 h-24"
                        value={description}
                        onChange={e => setDescription(e.target.value)}
                        required
                    />
                </div>

                <div className="grid grid-cols-2 gap-4">
                    <div>
                        <label className="block text-sm font-medium mb-1">Priority</label>
                        <select
                            className="w-full border rounded p-2"
                            value={priority}
                            onChange={e => setPriority(e.target.value)}
                        >
                            <option value="Low">Low</option>
                            <option value="Standard">Standard</option>
                            <option value="Critical">Critical</option>
                        </select>
                    </div>
                </div>

                <button className="bg-blue-600 text-white px-6 py-2 rounded font-medium hover:bg-blue-700">
                    Create Ticket
                </button>
            </form>
        </div>
    );
}
