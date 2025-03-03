document.addEventListener('DOMContentLoaded', function() {
    // Get DOM elements
    const giveBtn = document.getElementById('giveBtn');
    const takebackBtn = document.getElementById('takebackBtn');
    const newFriendBtn = document.getElementById('newFriendBtn');
    
    const giveSection = document.getElementById('giveSection');
    const takebackSection = document.getElementById('takebackSection');
    const newFriendSection = document.getElementById('newFriendSection');
    
    const friendListGive = document.getElementById('friendListGive');
    const friendListTakeback = document.getElementById('friendListTakeback');
    const itemListTakeback = document.getElementById('itemListTakeback');
    
    const selectedFriendGive = document.getElementById('selectedFriendGive');
    const itemNameGive = document.getElementById('itemNameGive');
    const submitGive = document.getElementById('submitGive');
    const giveResult = document.getElementById('giveResult');
    
    const selectedFriendTakeback = document.getElementById('selectedFriendTakeback');
    const selectedItemTakeback = document.getElementById('selectedItemTakeback');
    const submitTakeback = document.getElementById('submitTakeback');
    const takebackResult = document.getElementById('takebackResult');
    
    const newFriendName = document.getElementById('newFriendName');
    const submitNewFriend = document.getElementById('submitNewFriend');
    const newFriendResult = document.getElementById('newFriendResult');
    
    // Global state variables
    let selectedFriendId = null;
    let selectedItemId = null;
    
    // API endpoints (adjust these to match your Go backend)
    const API_URL = 'http://localhost:8080/api'; // Adjust this to your backend URL

    // Create a toast notification container
    const toastContainer = document.createElement('div');
    toastContainer.id = 'toast-container';
    document.body.appendChild(toastContainer);

    // Function to show a toast notification
    function showToast(message, type = 'info') {
        const toast = document.createElement('div');
        toast.className = `toast-notification ${type}`;
        
        // Create icon element
        const iconElement = document.createElement('div');
        iconElement.className = 'toast-icon';
        
        // Set icon content based on type
        let iconContent = '✓';
        if (type === 'error') iconContent = '!';
        if (type === 'info') iconContent = 'i';
        
        iconElement.textContent = iconContent;
        
        // Create message element
        const messageElement = document.createElement('div');
        messageElement.textContent = message;
        
        // Add elements to toast
        toast.appendChild(iconElement);
        toast.appendChild(messageElement);
        
        // Add to container
        toastContainer.appendChild(toast);
        
        // Show the toast
        setTimeout(() => {
            toast.classList.add('show');
        }, 10);
        
        // Hide and remove the toast after 5 seconds
        setTimeout(() => {
            toast.classList.remove('show');
            setTimeout(() => {
                toastContainer.removeChild(toast);
            }, 500);
        }, 5000);
    }

    // Show the selected section and hide others
    function showSection(section) {
        giveSection.style.display = 'none';
        takebackSection.style.display = 'none';
        newFriendSection.style.display = 'none';
        
        section.style.display = 'block';
        
        // Reset selections when switching sections
        resetSelections();
    }
    
    // Reset all selections
    function resetSelections() {
        selectedFriendId = null;
        selectedItemId = null;
        selectedFriendGive.value = '';
        selectedFriendTakeback.value = '';
        selectedItemTakeback.value = '';
        itemNameGive.value = '';
        newFriendName.value = '';
        
        hideMessages();
    }
    
    // Hide all result messages
    function hideMessages() {
        giveResult.style.display = 'none';
        takebackResult.style.display = 'none';
        newFriendResult.style.display = 'none';
    }
    
    // Fetch all friends from the API
    async function fetchFriends() {
        try {
            const response = await fetch(`${API_URL}/friends`);
            if (!response.ok) {
                throw new Error('Failed to fetch friends');
            }
            const friends = await response.json();
            return friends;
        } catch (error) {
            console.error('Error fetching friends:', error);
            showToast('Failed to load your friends. Please try again.', 'error');
            return [];
        }
    }
    
    // Display friends in the UI
    function displayFriends(containerId, friends) {
        const container = document.getElementById(containerId);
        
        if (friends.length === 0) {
            container.innerHTML = '<p>No friends found. Add a friend first.</p>';
            return;
        }
        
        container.innerHTML = '';
        friends.forEach(friend => {
            const friendElement = document.createElement('div');
            friendElement.classList.add('friend-item');
            friendElement.textContent = friend.name;
            friendElement.dataset.id = friend.id;
            friendElement.dataset.name = friend.name;
            
            container.appendChild(friendElement);
        });
    }
    
    // Fetch items for a specific friend
    async function fetchItems(friendId) {
        try {
            const response = await fetch(`${API_URL}/friends/${friendId}/items`);
            if (!response.ok) {
                throw new Error('Failed to fetch items');
            }
            const items = await response.json();
            return items;
        } catch (error) {
            console.error('Error fetching items:', error);
            showToast('Failed to load items from this friend.', 'error');
            return [];
        }
    }
    
    // Display items in the UI
    function displayItems(items) {
        if (items.length === 0) {
            itemListTakeback.innerHTML = '<p>No items found for this friend.</p>';
            return;
        }
        
        itemListTakeback.innerHTML = '';
        items.forEach(item => {
            const itemElement = document.createElement('div');
            itemElement.classList.add('item-item');
            itemElement.textContent = item.name;
            itemElement.dataset.id = item.id;
            itemElement.dataset.name = item.name;
            
            itemListTakeback.appendChild(itemElement);
        });
    }
    
    // Initialize the app
    async function initApp() {
        const friends = await fetchFriends();
        displayFriends('friendListGive', friends);
        displayFriends('friendListTakeback', friends);
    }
    
    // Event listeners for action buttons
    giveBtn.addEventListener('click', function() {
        showSection(giveSection);
        initApp();
    });
    
    takebackBtn.addEventListener('click', function() {
        showSection(takebackSection);
        initApp();
    });
    
    newFriendBtn.addEventListener('click', function() {
        showSection(newFriendSection);
    });
    
    // Event delegation for friend selection in Give section
    friendListGive.addEventListener('click', function(e) {
        if (e.target.classList.contains('friend-item')) {
            // Remove selected class from all friends
            const friendItems = friendListGive.querySelectorAll('.friend-item');
            friendItems.forEach(item => item.classList.remove('selected'));
            
            // Add selected class to clicked friend
            e.target.classList.add('selected');
            
            // Update selected friend
            selectedFriendId = parseInt(e.target.dataset.id, 10);
            selectedFriendGive.value = e.target.dataset.name;
            
            // Show acknowledgment
            showToast(`Selected ${e.target.dataset.name} to receive an item`, 'info');
        }
    });
    
    // Event delegation for friend selection in Takeback section
    friendListTakeback.addEventListener('click', async function(e) {
        if (e.target.classList.contains('friend-item')) {
            // Remove selected class from all friends
            const friendItems = friendListTakeback.querySelectorAll('.friend-item');
            friendItems.forEach(item => item.classList.remove('selected'));
            
            // Add selected class to clicked friend
            e.target.classList.add('selected');
            
            // Update selected friend
            selectedFriendId = parseInt(e.target.dataset.id, 10);
            selectedFriendTakeback.value = e.target.dataset.name;
            
            // Show acknowledgment
            showToast(`Looking at items borrowed by ${e.target.dataset.name}`, 'info');
            
            // Fetch and display items for the selected friend
            const items = await fetchItems(selectedFriendId);
            displayItems(items);
        }
    });
    
    // Event delegation for item selection in Takeback section
    itemListTakeback.addEventListener('click', function(e) {
        if (e.target.classList.contains('item-item')) {
            // Remove selected class from all items
            const itemItems = itemListTakeback.querySelectorAll('.item-item');
            itemItems.forEach(item => item.classList.remove('selected'));
            
            // Add selected class to clicked item
            e.target.classList.add('selected');
            
            // Update selected item
            selectedItemId = e.target.dataset.id;
            selectedItemTakeback.value = e.target.dataset.name;
            
            // Show acknowledgment
            showToast(`Selected "${e.target.dataset.name}" to take back`, 'info');
        }
    });
    
    // Give item form submission
    submitGive.addEventListener('click', async function() {
        hideMessages();
        
        // Validate inputs
        if (!selectedFriendId) {
            giveResult.textContent = 'Please select a friend.';
            giveResult.style.display = 'block';
            return;
        }
        
        if (!itemNameGive.value.trim()) {
            giveResult.textContent = 'Please enter an item name.';
            giveResult.style.display = 'block';
            return;
        }
        
        // Submit the form
        try {
            const response = await fetch(`${API_URL}/items`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    friendId: selectedFriendId,
                    name: itemNameGive.value.trim()
                }),
            });
            
            if (!response.ok) {
                 // Try to parse JSON error response
                try {
                    const errorData = await response.json();
                    throw new Error(`Failed to add item: ${errorData.error || 'Unknown error'}`);
                } catch (parseError) {
                    // If JSON parsing fails, use the default error message
                    throw new Error('Failed to add item');
                }
            }
            
            const successMessage = `Successfully lent ${itemNameGive.value.trim()} to ${selectedFriendGive.value}.`;
            giveResult.textContent = successMessage;
            giveResult.style.display = 'block';
            giveResult.className = 'result acknowledgment';
            
            // Show toast notification with acknowledgment
            showToast(successMessage, 'success');
            
            // Reset form
            itemNameGive.value = '';
            resetSelections();
            
        } catch (error) {
            console.error('Error adding item:', error);
            giveResult.textContent = error.message || 'An error occurred. Please try again.';
            giveResult.style.display = 'block';
            
            // Show toast notification with error
            showToast(error.message || 'Failed to lend the item. Please try again.', 'error');
        }
    });
    
    // Takeback item form submission
    submitTakeback.addEventListener('click', async function() {
        hideMessages();
        
        // Validate inputs
        if (!selectedFriendId) {
            takebackResult.textContent = 'Please select a friend.';
            takebackResult.style.display = 'block';
            return;
        }
        
        if (!selectedItemId) {
            takebackResult.textContent = 'Please select an item.';
            takebackResult.style.display = 'block';
            return;
        }
        
        // Submit the form
        try {
            const response = await fetch(`${API_URL}/items/${selectedItemId}`, {
                method: 'DELETE',
            });
            
            if (!response.ok) {
                throw new Error('Failed to delete item');
            }
            
            const successMessage = `Successfully took back ${selectedItemTakeback.value} from ${selectedFriendTakeback.value}.`;
            takebackResult.textContent = successMessage;
            takebackResult.style.display = 'block';
            takebackResult.className = 'result acknowledgment';
            
            // Show toast notification with acknowledgment
            showToast(successMessage, 'success');
            
            // Refresh items list
            const items = await fetchItems(selectedFriendId);
            displayItems(items);
            
            // Reset selection
            selectedItemId = null;
            selectedItemTakeback.value = '';
            resetSelections();
            
        } catch (error) {
            console.error('Error deleting item:', error);
            takebackResult.textContent = 'An error occurred. Please try again.';
            takebackResult.style.display = 'block';
            
            // Show toast notification with error
            showToast('Failed to take back the item. Please try again.', 'error');
        }
    });
    
    // New friend form submission
    submitNewFriend.addEventListener('click', async function() {
        hideMessages();
        
        // Validate input
        if (!newFriendName.value.trim()) {
            newFriendResult.textContent = 'Please enter a friend name.';
            newFriendResult.style.display = 'block';
            return;
        }
        
        // Submit the form
        try {
            const response = await fetch(`${API_URL}/friends`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    name: newFriendName.value.trim()
                }),
            });
            
            if (!response.ok) {
                throw new Error('Failed to add friend');
            }
            
            const successMessage = `Successfully added ${newFriendName.value.trim()} as a friend.`;
            newFriendResult.textContent = successMessage;
            newFriendResult.style.display = 'block';
            newFriendResult.className = 'result acknowledgment';
            
            // Show toast notification with acknowledgment
            showToast(successMessage, 'success');
            
            // Reset form
            newFriendName.value = '';
            initApp();
            resetSelections();
            showSection(newFriendSection);
            
        } catch (error) {
            console.error('Error adding friend:', error);
            newFriendResult.textContent = 'An error occurred. Please try again.';
            newFriendResult.style.display = 'block';
            
            // Show toast notification with error
            showToast('Failed to add friend. Please try again.', 'error');
        }
    });
    
    // Show the new friend section by default
    showSection(newFriendSection);
});