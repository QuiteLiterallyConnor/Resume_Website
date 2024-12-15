$(document).ready(function() {
    const userId = window.location.pathname.split("/").pop();

    function fetchUserDetails() {
        $.getJSON(`/api/users/${userId}`, function(user) {
            const userDetailsDiv = $("#userDetails");

            userDetailsDiv.empty(); // Clear any existing data

            // Ensure AccessedParts is properly processed
            const accessedPartsList = user.Accessed_Parts && user.Accessed_Parts.length > 0 ?
                user.Accessed_Parts.map(part => `
                    <div class="row mb-2">
                        <div class="col-6"><strong>Part Accessed:</strong> ${part.part}</div>
                        <div class="col-6"><strong>Time Accessed:</strong> ${part.time_accessed}</div>
                    </div>
                `).join('') : 
                '<p>No parts accessed</p>';

            const userInfo = `
                <p><strong>ID:</strong> ${user.ID}</p>
                <p><strong>IP Address:</strong> ${user.IP_Address}</p>
                <p><strong>Country:</strong> ${user.Country.String}</p>
                <p><strong>City:</strong> ${user.City.String}</p>
                <p><strong>First Time Accessed:</strong> ${user.First_Time_Accessed}</p>
                <p><strong>Last Time Accessed:</strong> ${user.Last_Time_Accessed}</p>
                <p><strong>Blacklisted:</strong> ${user.Blacklisted ? "Yes" : "No"}</p>
                <p><strong>Client Data:</strong> ${user.Client_Data}</p>
                <h3>Accessed Parts</h3>
                <div class="container">${accessedPartsList}</div>
                <div class="mt-3">
                    <button class="btn btn-danger" id="deleteButton">Delete User</button>
                    <button class="btn btn-warning ml-2" id="toggleBlacklistButton">${user.Blacklisted ? "Unblacklist" : "Blacklist"}</button>
                </div>
            `;

            userDetailsDiv.append(userInfo);

            $("#deleteButton").click(function() {
                deleteUser(user.ID);
            });

            $("#toggleBlacklistButton").click(function() {
                toggleBlacklist(user.ID, !user.Blacklisted);
            });
        }).fail(function() {
            console.error("Error fetching user details.");
        });
    }

    function deleteUser(userId) {
        $.ajax({
            url: `/api/users/${userId}`,
            type: 'DELETE',
            success: function(result) {
                window.location.href = "/";
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
                fetchUserDetails(); // Refresh the user details after updating blacklist status
            },
            error: function() {
                console.error("Error updating blacklist status.");
            }
        });
    }

    // Back button to return to the main user list
    $("#backButton").click(function() {
        window.location.href = "/";
    });

    // Initial fetch of user details
    fetchUserDetails();
});
