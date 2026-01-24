import CreateUserCard from "../components/Users/CreateUserCard";
import UsersTable from "../components/Users/UsersTable";

export default function AdminUsers() {
  return (
    <div className="w-full lg:mx-auto lg:max-w-7xl">
      <h1 className="mb-6 px-4 sm:px-6 lg:px-8 text-2xl font-semibold text-slate-900">
        Users
      </h1>

      <div className="px-0 sm:px-0 lg:px-0">
        <CreateUserCard />
        <UsersTable />
      </div>
    </div>
  );
}


