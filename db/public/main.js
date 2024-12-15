$(document).ready(function() {
    let usersData = [];
    let filteredData = [];
    let sortOrder = {
        ID: 'asc',
        IP_Address: 'asc',
        Country: 'asc',
        City: 'asc',
        First_Time_Accessed: 'asc',
        Last_Time_Accessed: 'desc', // Set the initial sort order for "Last Time Accessed" to descending
        Blacklisted: 'asc'
    };

    function fetchUsers() {
        $.getJSON("/api/users", function(data) {
            usersData = data; // Store the data for sorting and filtering
            filteredData = [...usersData]; // Initialize filtered data with the full data set
            
            // Initially sort by "Last Time Accessed" in descending order
            sortData('Last_Time_Accessed');
        }).fail(function() {
            console.error("Error fetching user data.");
        });
    }

    function renderTable(data) {
        const tableBody = $("#userTable tbody");
        tableBody.empty(); // Clear the table before appending new data
        data.forEach(function(user) {
            const row = $("<tr></tr>");

            const idCell = $("<td></td>").text(user.ID);
            row.append(idCell);

            const ipCell = $("<td></td>").text(user.IP_Address);
            row.append(ipCell);

            const countryCell = $("<td></td>").text(user.Country.String || "Unknown");
            row.append(countryCell);

            const cityCell = $("<td></td>").text(user.City.String || "Unknown");
            row.append(cityCell);

            const firstTimeCell = $("<td></td>").text(user.First_Time_Accessed);
            row.append(firstTimeCell);

            const lastTimeCell = $("<td></td>").text(user.Last_Time_Accessed);
            row.append(lastTimeCell);

            const blacklistedCell = $("<td></td>").text(user.Blacklisted ? "Yes" : "No");
            row.append(blacklistedCell);

            const actionsCell = $("<td></td>");

            const deleteButton = $("<button></button>")
                .addClass("btn btn-danger btn-sm")
                .text("Delete")
                .click(function() {
                    deleteUser(user.ID);
                });

            const toggleBlacklistButton = $("<button></button>")
                .addClass("btn btn-warning btn-sm ml-2")
                .text(user.Blacklisted ? "Unblacklist" : "Blacklist")
                .click(function() {
                    toggleBlacklist(user.ID, !user.Blacklisted);
                });

            const moreInfoButton = $("<button></button>")
                .addClass("btn btn-info btn-sm ml-2")
                .text("More Info")
                .click(function() {
                    window.location.href = `/user/${user.ID}`;
                });

            actionsCell.append(deleteButton, toggleBlacklistButton, moreInfoButton);
            row.append(actionsCell);

            tableBody.append(row);
        });
    }

    function sortData(field) {
        const order = sortOrder[field];
        filteredData.sort((a, b) => {
            let valA = a[field];
            let valB = b[field];

            // Handle null or undefined values
            if (valA == null) valA = "";
            if (valB == null) valB = "";

            // Convert to lowercase if string to ensure case-insensitive sorting
            if (typeof valA === 'string') valA = valA.toLowerCase();
            if (typeof valB === 'string') valB = valB.toLowerCase();

            // Handle sorting for Country and City (checking for null values within the object)
            if (field === 'Country' || field === 'City') {
                valA = a[field].String || "Unknown";
                valB = b[field].String || "Unknown";
            }

            // Handle date comparisons for First_Time_Accessed and Last_Time_Accessed
            if (field === 'First_Time_Accessed' || field === 'Last_Time_Accessed') {
                valA = new Date(valA).getTime() || 0;
                valB = new Date(valB).getTime() || 0;
            }

            // Handle boolean comparison for Blacklisted
            if (field === 'Blacklisted') {
                valA = valA ? 1 : 0;  // Convert true to 1 and false to 0
                valB = valB ? 1 : 0;  // Convert true to 1 and false to 0
            }

            // Handle numeric comparison
            if (!isNaN(valA) && !isNaN(valB)) {
                valA = parseFloat(valA);
                valB = parseFloat(valB);
            }

            // Compare values
            if (valA < valB) return order === 'asc' ? -1 : 1;
            if (valA > valB) return order === 'asc' ? 1 : -1;
            return 0;
        });
        sortOrder[field] = order === 'asc' ? 'desc' : 'asc'; // Toggle sort order
        updateSortIcons(field, sortOrder[field]); // Update the sort icon based on the current order
        renderTable(filteredData);
    }

    function updateSortIcons(activeField, activeOrder) {
        // Reset all icons to the default (unsorted) state
        $(".sort-icon").removeClass("bi-sort-alpha-down bi-sort-alpha-up").addClass("bi-sort-alpha-down");

        // Update the icon for the active field
        const iconClass = activeOrder === 'asc' ? "bi-sort-alpha-down" : "bi-sort-alpha-up";
        $(`.sort-icon[data-sort="${activeField}"]`).removeClass("bi-sort-alpha-down").addClass(iconClass);
    }

    function filterData() {
        const searchIP = $("#searchIP").val().toLowerCase();
        const searchCountry = $("#searchCountry").val().toLowerCase();
        const searchCity = $("#searchCity").val().toLowerCase();

        filteredData = usersData.filter(user => {
            const ipMatch = user.IP_Address.toLowerCase().includes(searchIP);
            const countryMatch = (user.Country.String || "Unknown").toLowerCase().includes(searchCountry);
            const cityMatch = (user.City.String || "Unknown").toLowerCase().includes(searchCity);
            return ipMatch && countryMatch && cityMatch;
        });

        renderTable(filteredData);
    }

    function deleteUser(userId) {
        $.ajax({
            url: `/api/users/${userId}`,
            type: 'DELETE',
            success: function(result) {
                fetchUsers(); // Refresh the user list after deletion
            },
            error: function() {
                console.error("Error deleting user.");
            }
        });
    }

    function toggleBlacklist(userId, blacklistStatus) {
        $.ajax({
            url: `/api/users/${userId}/blacklist`,
            type: 'PATCH',
            data: JSON.stringify({ blacklisted: blacklistStatus }),
            contentType: 'application/json',
            success: function(result) {
                fetchUsers(); // Refresh the user list after updating blacklist status
            },
            error: function() {
                console.error("Error updating blacklist status.");
            }
        });
    }

    // Event listeners for sort icons
    $('.sort-icon').click(function() {
        const field = $(this).data('sort');
        sortData(field);
    });

    // Event listeners for search fields
    $("#searchIP, #searchCountry, #searchCity").on("input", function() {
        filterData();
    });

    // Initial fetch of users
    fetchUsers();
});
