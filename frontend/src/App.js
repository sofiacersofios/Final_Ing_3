// app.js
import React, { useState, useEffect } from 'react';
import './App.css';

function App() {
	const [data, setData] = useState([]);
	const [newItem, setNewItem] = useState({ name: '' });

	useEffect(() => {
		fetchData();
	}, []);

	const fetchData = async () => {
		const response = await fetch('http://localhost:8080/api/data');
		const result = await response.json();
		setData(result);
	};

	const handleInputChange = (e) => {
		setNewItem({ ...newItem, [e.target.name]: e.target.value });
	};

	const handleAddItem = async () => {
		const response = await fetch('http://localhost:8080/api/data', {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json',
			},
			body: JSON.stringify(newItem),
		});

		if (response.status === 201) {
			// Item created successfully
			fetchData(); // Refresh the data
			setNewItem({ name: '' }); // Clear the input
		} else {
			// Handle error
			console.error('Failed to add item');
		}
	};

	const handleDeleteItem = async (itemId) => {
		try {
			const response = await fetch(`http://localhost:8080/api/data/${itemId}`, {
				method: 'DELETE',
			});
	
			if (response.status === 200) {
				// Item deleted successfully
				fetchData(); // Refresh the data
			} else {
				const errorText = await response.text();
				console.error('Failed to delete item:', response.statusText, errorText);
			}
		} catch (error) {
			console.error('Error during DELETE request:', error.message);
		}
	};
	

	return (
		<div>
			<h1>Data from Backend</h1>
			<ul>
				{data.map((item) => (
					<li key={item.id}>
						{item.name}
						<button className='x' onClick={() => handleDeleteItem(item.id)}>X</button>
					</li>
				))}
			</ul>
			<div>
				<input
					type="text"
					name="name"
					value={newItem.name}
					onChange={handleInputChange}
					placeholder="Enter item name"
				/>
				<button onClick={handleAddItem}>Add Item</button>
			</div>
		</div>
	);
}

export default App;
