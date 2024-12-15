document.addEventListener('DOMContentLoaded', function() {
    const options = document.querySelectorAll('.option');
    const segments = document.querySelectorAll('.segment-box');
    // const droneImage = document.getElementById('drone-image');

    // $('.segment-box').on('click', function() {
    //     const imageSrc = $(this).data('image');
    //     const droneImage = $('#drone-image');
    //     if (droneImage.length) {
    //         droneImage.attr('src', imageSrc);
    //     } else {
    //         console.error('Element #drone-image not found');
    //     }
    // });

    options.forEach(option => {
        option.addEventListener('mouseenter', () => {
            const newImage = option.getAttribute('data-image');
            droneImage.src = newImage;
            options.forEach(opt => opt.classList.remove('active'));
            option.classList.add('active');
        });
    });

    segments.forEach(segment => {
        segment.addEventListener('mouseenter', () => {
            const newImage = segment.getAttribute('data-image');
            droneImage.src = newImage;
        });
    });
});
