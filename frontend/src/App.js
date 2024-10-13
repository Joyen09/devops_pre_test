import React, { useState, useEffect, useMemo } from 'react';
import axios from 'axios';
import { useTable, usePagination } from 'react-table';
import { ToastContainer, toast } from 'react-toastify';
import 'react-toastify/dist/ReactToastify.css';
import { API_BASE_URL } from './config';

function App() {
    const [customers, setCustomers] = useState([]);
    const [selectedCustomer, setSelectedCustomer] = useState(null);
    const [transactions, setTransactions] = useState([]);
    const [newCustomer, setNewCustomer] = useState({ name: '', email: '' });
    const [startDate, setStartDate] = useState('');
    const [endDate, setEndDate] = useState('');

    useEffect(() => {
        fetchCustomers();
    }, []);

    const fetchCustomers = async () => {
        try {
            const response = await axios.get(`${API_BASE_URL}/api/customers`);
            setCustomers(response.data);
        } catch (error) {
            console.error('Error fetching customers:', error);
            toast.error('獲取客戶清單失敗');
        }
    };

    const fetchCustomerDetails = async (id) => {
        try {
            const response = await axios.get(`${API_BASE_URL}/api/customers/${id}`);
            setSelectedCustomer(response.data);
        } catch (error) {
            console.error('Error fetching customer details:', error);
            toast.error('獲取客戶詳情失敗');
        }
    };

    const fetchCustomerTransactions = async (id) => {
        try {
            if (!startDate || !endDate) {
                toast.warning('請選擇開始和結束日期');
                return;
            }
            const response = await axios.get(`${API_BASE_URL}/api/customers/${id}/transactions`, {
                params: { start_date: startDate, end_date: endDate },
            });
            if (response.data === null) {
                console.log('No transactions found or error occurred');
                toast.info('在所選日期範圍內沒有找到交易記錄，或發生了錯誤');
                setTransactions([]);
            } else if (Array.isArray(response.data)) {
                if (response.data.length === 0) {
                    toast.info('所選日期範圍內沒有交易記錄');
                }
                setTransactions(response.data);
            } else {
                console.error('Unexpected response format:', response.data);
                toast.error('獲取交易記錄失敗：數據格式錯誤');
                setTransactions([]);
            }
        } catch (error) {
            console.error('Error fetching transactions:', error);
            toast.error('獲取交易記錄失敗');
            setTransactions([]);
        }
    };

    const handleCreateCustomer = async () => {
        try {
            await axios.post(`${API_BASE_URL}/api/customers`, newCustomer);
            setNewCustomer({ name: '', email: '' });
            fetchCustomers();
            toast.success('客戶創建成功');
        } catch (error) {
            console.error('Error creating customer:', error);
            toast.error('創建客戶失敗');
        }
    };

    const handleUpdateCustomer = async () => {
        try {
            await axios.put(`${API_BASE_URL}/api/customers/${selectedCustomer.customer.id}`, {
                name: selectedCustomer.customer.name,
                email: selectedCustomer.customer.email,
            });
            fetchCustomers();
            toast.success('客戶資訊更新成功');
        } catch (error) {
            console.error('Error updating customer:', error);
            toast.error('更新客戶資訊失敗');
        }
    };

    const columns = useMemo(
        () => [
            {
                Header: '姓名',
                accessor: 'name',
            },
            {
                Header: '郵箱',
                accessor: 'email',
            },
        ],
        []
    );

    const {
        getTableProps,
        getTableBodyProps,
        headerGroups,
        prepareRow,
        page,
        canPreviousPage,
        canNextPage,
        pageOptions,
        pageCount,
        gotoPage,
        nextPage,
        previousPage,
        setPageSize,
        state: { pageIndex, pageSize },
    } = useTable(
        {
            columns,
            data: customers,
            initialState: { pageIndex: 0, pageSize: 10 },
        },
        usePagination
    );

    return (
        <div className="App">
            <ToastContainer />
            <h1>客戶管理系統</h1>

            <h2>客戶清單</h2>
            <table {...getTableProps()} style={{ border: 'solid 1px blue' }}>
                <thead>
                {headerGroups.map(headerGroup => (
                    <tr {...headerGroup.getHeaderGroupProps()}>
                        {headerGroup.headers.map(column => (
                            <th
                                {...column.getHeaderProps()}
                                style={{
                                    borderBottom: 'solid 3px red',
                                    background: 'aliceblue',
                                    color: 'black',
                                    fontWeight: 'bold',
                                }}
                            >
                                {column.render('Header')}
                            </th>
                        ))}
                    </tr>
                ))}
                </thead>
                <tbody {...getTableBodyProps()}>
                {page.map(row => {
                    prepareRow(row)
                    return (
                        <tr {...row.getRowProps()} onClick={() => fetchCustomerDetails(row.original.id)}>
                            {row.cells.map(cell => {
                                return (
                                    <td
                                        {...cell.getCellProps()}
                                        style={{
                                            padding: '10px',
                                            border: 'solid 1px gray',
                                            background: 'papayawhip',
                                        }}
                                    >
                                        {cell.render('Cell')}
                                    </td>
                                )
                            })}
                        </tr>
                    )
                })}
                </tbody>
            </table>
            <div className="pagination">
                <button onClick={() => gotoPage(0)} disabled={!canPreviousPage}>
                    {'<<'}
                </button>{' '}
                <button onClick={() => previousPage()} disabled={!canPreviousPage}>
                    {'<'}
                </button>{' '}
                <button onClick={() => nextPage()} disabled={!canNextPage}>
                    {'>'}
                </button>{' '}
                <button onClick={() => gotoPage(pageCount - 1)} disabled={!canNextPage}>
                    {'>>'}
                </button>{' '}
                <span>
                    Page{' '}
                    <strong>
                        {pageIndex + 1} of {pageOptions.length}
                    </strong>{' '}
                </span>
                <span>
                    | Go to page:{' '}
                    <input
                        type="number"
                        defaultValue={pageIndex + 1}
                        onChange={e => {
                            const page = e.target.value ? Number(e.target.value) - 1 : 0
                            gotoPage(page)
                        }}
                        style={{ width: '50px' }}
                    />
                </span>{' '}
                <select
                    value={pageSize}
                    onChange={e => {
                        setPageSize(Number(e.target.value))
                    }}
                >
                    {[10, 20, 30, 40, 50].map(pageSize => (
                        <option key={pageSize} value={pageSize}>
                            Show {pageSize}
                        </option>
                    ))}
                </select>
            </div>

            {selectedCustomer && (
                <div>
                    <h2>客戶詳情</h2>
                    <input
                        value={selectedCustomer.customer.name}
                        onChange={(e) => setSelectedCustomer({
                            ...selectedCustomer,
                            customer: {...selectedCustomer.customer, name: e.target.value}
                        })}
                    />
                    <input
                        value={selectedCustomer.customer.email}
                        onChange={(e) => setSelectedCustomer({
                            ...selectedCustomer,
                            customer: {...selectedCustomer.customer, email: e.target.value}
                        })}
                    />
                    <button onClick={handleUpdateCustomer}>更新客戶資料</button>
                    <p>過去一年交易總額: {selectedCustomer.total_amount_last_year}</p>

                    <h3>交易記錄 - {selectedCustomer.customer.name} (ID: {selectedCustomer.customer.id})</h3>
                    <input
                        type="date"
                        value={startDate}
                        onChange={(e) => setStartDate(e.target.value)}
                    />
                    <input
                        type="date"
                        value={endDate}
                        onChange={(e) => setEndDate(e.target.value)}
                    />
                    <button onClick={() => fetchCustomerTransactions(selectedCustomer.customer.id)}>
                        查詢交易
                    </button>
                    {transactions && transactions.length > 0 ? (
                        <ul>
                            {transactions.map(transaction => (
                                <li key={transaction.id}>
                                    日期：{new Date(transaction.transaction_date).toLocaleDateString()} -
                                    金額：{transaction.borrow_fee}
                                </li>
                            ))}
                        </ul>
                    ) : (
                        <p>目前沒有可顯示的交易記錄</p>
                    )}
                </div>
            )}

            <h2>新增客戶</h2>
            <input
                placeholder="姓名"
                value={newCustomer.name}
                onChange={(e) => setNewCustomer({...newCustomer, name: e.target.value})}
            />
            <input
                placeholder="電子郵件"
                value={newCustomer.email}
                onChange={(e) => setNewCustomer({...newCustomer, email: e.target.value})}
            />
            <button onClick={handleCreateCustomer}>新增客戶</button>
        </div>
    );
}

export default App;
