package psql

const (
	categoryIsExist = iota

	storeInventoriesQ
	storeGetStatusQ
	storeCreateQ
	storeFindByIDQ
	storeDeleteByIDQ
	storeCreateInvoiceByDatesQ

	userCreateQ
	userGetAllowedMethodsAndPassQ
	userGetByNameQ
	userUpdateQ
	userDeleteQ

	petGetStatusIDByNameQ
	petAddToStoreQ
	petUpdateWithBodyQ
	petFindByStatusQ
	petGetByIDQ
	petUpdateWithFieldsQ
	petDeleteByIDQ
	petUpdatePhotosByIDQ

	tagsInsertQ
	tagsDeleteQ
	tagsGetByPetIDQ
)

// query master
var qm = map[int]string{
	categoryIsExist: `
	SELECT count(*) FROM information_schema.TABLES
	WHERE (TABLE_SCHEMA = 'public') AND (TABLE_NAME = 'category');`,

	storeInventoriesQ: `
	select pet_status, sum(quantity) from order_info group by pet_status`,

	storeGetStatusQ: `
	select id from order_status where name=$1`,

	storeCreateQ: `
	insert into "order" (pet_id, user_id, quantity, ship_date, order_status_id, complete)
	values (:pet_id, :user_id, :quantity, :ship_date, :order_status, :complete) returning id;`,

	storeFindByIDQ: `
	select id, user_id, pet_id, quantity,
	ship_date, order_status, complete 
	from order_info where id=$1`,

	storeDeleteByIDQ: `
	delete from "order" 
	where id=$1 returning id`,

	storeCreateInvoiceByDatesQ: `
	select id, user_name, pet, category, ship_date, quantity, price
	from invoice_info where ship_date between :from and :to
	group by category, ship_date, id, user_name, pet, quantity, price
	order by category, ship_date;
	`,

	userCreateQ: `
	insert into "user" (user_name, first_name, 
	last_name, email, password, phone, user_status_id) 
	values (:user_name, :first_name, :last_name, 
	:email, :password, :phone, :user_status_id)`,

	userGetAllowedMethodsAndPassQ: `
	select password, allowed_methods 
	from user_info where user_name=$1`,

	userGetByNameQ: `
	select id, user_name, password, first_name,
	last_name, email, phone, user_status_id
	from "user" where user_name = $1`,

	userUpdateQ: `update "user" set 
	user_name = :user_name, first_name = :first_name, last_name = :last_name, 
	email = :email, password = :password, phone = :phone, user_status_id = :user_status_id 
	where user_name = :user_name`,

	userDeleteQ: `delete from "user" where user_name=$1`,

	petGetStatusIDByNameQ: `
	select id from pet_status where name = $1`,

	petAddToStoreQ: `
	insert into pet 
    (category_id, name, photo_urls, pet_status_id)
 	values (:category_id, :name, :photo_urls, :pet_status_id) returning id;`,

	petUpdateWithBodyQ: `
	update pet set category_id = :category_id, name = :name, 
	photo_urls = :photo_urls, pet_status_id = :pet_status_id
	where id = :id returning id;`,

	petFindByStatusQ: `
	select id, category_id, category_name, name,
	photo_urls, pet_status_name from pet_info
	where pet_status_name=any($1)`,

	petGetByIDQ: `
	select id, category_id, 
	category_name, name, photo_urls,
	pet_status_name from pet_info where id=$1`,

	petUpdateWithFieldsQ: `
	update pet set name=:name, pet_status_id=:pet_status_id 
	where id = :id`,

	petDeleteByIDQ: `
	delete from pet 
	where id=$1`,

	petUpdatePhotosByIDQ: `
	update pet set photo_urls=:photo_urls 
	where id=:id`,

	tagsInsertQ: `
	insert into pet_tag 
	VALUES (:pet_id, :tag_id)`,

	tagsDeleteQ: `
	delete from pet_tag 
	where pet_id = $1`,

	tagsGetByPetIDQ: `
	select t.id, t.name from pet_tag 
	inner join tag t on pet_tag.tag_id = t.id 
	where pet_id=$1`,
}
