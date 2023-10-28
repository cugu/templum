
document.addEventListener('DOMContentLoaded', () => {
    const toggles = document.querySelectorAll('.toggle');
    // Add a click listener to each button
    toggles.forEach(question =>
        question.addEventListener('click', ({ target }) => {
            target.parentElement.parentElement.classList.toggle('is-active');

            // no default action
            return false;
        })
    );
})
