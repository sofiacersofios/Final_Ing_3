const path = require('path');

Feature('App Functionality');

Scenario('should render the app ...', ({ I }) => {
  I.amOnPage('/');
  I.see('NOTAS');
  I.seeElement('ul');
  I.seeElement('input[name="name"]');
  I.seeElement('button');
  I.wait(3);

});

Scenario('should add a new item', async ({ I }) => {
  I.amOnPage('/');

  I.wait(3);

  // Get the initial count of items
  const initialItemCount = await I.grabNumberOfVisibleElements('li');

  // Add a new item
  I.fillField('input[name="name"]', 'New Note');
  I.click('Add Note');

  // Wait for the new item to be added dynamically
  await I.waitForFunction(
    (initialCount) => {
      const currentCount = document.querySelectorAll('li').length;
      return currentCount === initialCount + 1;
    },
    [initialItemCount],
    30 // Set a reasonable timeout, e.g., 30 seconds
  );

  // Check if the item count increased using native JavaScript assertion
  const newItemCount = await I.grabNumberOfVisibleElements('li');
  if (newItemCount !== initialItemCount + 1) {
    throw new Error('Item count did not increase as expected');
  }

  I.wait(3);

});


Scenario('should delete the last added item', async ({ I }) => {
  I.amOnPage('/');

  I.wait(3);

  // Grab the count of items before deletion
  const originalItemCount = await I.grabNumberOfVisibleElements('ul li');

  // Click the delete button of the last item
  I.click('ul li:last-child button.x');

  // Wait for the item to be removed dynamically
  await I.waitForFunction(
    (originalItemCount) => {
      const items = document.querySelectorAll('ul li');
      return items.length < originalItemCount;
    },
    [originalItemCount],
    20
  );

  I.wait(3);
});



